package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type adminButton struct {
	Text string `json:"text"`
	Data string `json:"callback_data"`
}

type adminMarkup struct {
	Keyboard [][]adminButton `json:"inline_keyboard"`
}

type dockerContainer struct {
	ID     string            `json:"Id"`
	Names  []string          `json:"Names"`
	State  string            `json:"State"`
	Labels map[string]string `json:"Labels"`
}

func (b *Bot) sendAdminMenu(ctx context.Context, chatID int64) error {
	return b.sendMessageWithKeyboard(ctx, chatID, "Управление контейнерами проекта «"+b.config.ComposeProject+"». Действия затронут только контейнеры этого Compose-проекта.", adminMarkup{Keyboard: [][]adminButton{
		{{Text: "🔄 Перезапустить все", Data: "admin:restart"}},
		{{Text: "🗑 Удалить все контейнеры…", Data: "admin:confirm-delete"}},
	}})
}

func (b *Bot) handleAdminCallback(ctx context.Context, query *CallbackQuery) error {
	if query.Message == nil || b.config.AdminTelegramID == 0 || query.From.ID != b.config.AdminTelegramID || query.Message.Chat.Type != "private" || query.Message.Chat.ID != b.config.AdminTelegramID {
		return b.answerCallback(ctx, query.ID, "Нет доступа")
	}
	switch query.Data {
	case "admin:menu":
		_ = b.answerCallback(ctx, query.ID, "")
		return b.sendAdminMenu(ctx, query.Message.Chat.ID)
	case "admin:confirm-delete":
		_ = b.answerCallback(ctx, query.ID, "")
		return b.sendMessageWithKeyboard(ctx, query.Message.Chat.ID, "Подтвердите удаление всех контейнеров этого проекта. После удаления бот, база и Redis остановятся. Для возврата потребуется запустить docker compose up -d на сервере.", adminMarkup{Keyboard: [][]adminButton{
			{{Text: "Да, удалить контейнеры", Data: "admin:delete"}},
			{{Text: "Отмена", Data: "admin:menu"}},
		}})
	case "admin:restart":
		if err := b.answerCallback(ctx, query.ID, "Перезапускаю…"); err != nil {
			return err
		}
		count, err := b.manageProjectContainers(ctx, "restart")
		if err != nil {
			return b.sendMessage(ctx, query.Message.Chat.ID, "Не удалось перезапустить контейнеры: "+err.Error())
		}
		return b.sendMessage(ctx, query.Message.Chat.ID, fmt.Sprintf("Перезапущено контейнеров: %d.", count))
	case "admin:delete":
		if err := b.answerCallback(ctx, query.ID, "Удаляю контейнеры…"); err != nil {
			return err
		}
		count, err := b.manageProjectContainers(ctx, "delete")
		if err != nil {
			return b.sendMessage(ctx, query.Message.Chat.ID, "Удаление выполнено не полностью: "+err.Error())
		}
		return b.sendMessage(ctx, query.Message.Chat.ID, fmt.Sprintf("Удалено контейнеров: %d. Чтобы восстановить проект, выполните на сервере docker compose up -d.", count))
	default:
		return b.answerCallback(ctx, query.ID, "Кнопка устарела")
	}
}

func (b *Bot) manageProjectContainers(ctx context.Context, action string) (int, error) {
	if action != "restart" && action != "delete" {
		return 0, fmt.Errorf("unknown action")
	}
	client := b.dockerClient()
	containers, err := b.listProjectContainers(ctx, client)
	if err != nil {
		return 0, err
	}
	completed := 0
	for _, container := range containers {
		path := "/containers/" + url.PathEscape(container.ID)
		method := http.MethodPost
		if action == "restart" {
			path += "/restart?t=10"
		} else {
			method = http.MethodDelete
			path += "?force=true&v=false"
		}
		req, err := http.NewRequestWithContext(ctx, method, "http://docker"+path, nil)
		if err != nil {
			return completed, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return completed, err
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			return completed, fmt.Errorf("container %s: Docker API %s: %s", container.ID[:12], resp.Status, strings.TrimSpace(string(body)))
		}
		completed++
	}
	return completed, nil
}

func (b *Bot) restartProjectContainer(ctx context.Context, id string) error {
	containers, err := b.listProjectContainers(ctx, b.dockerClient())
	if err != nil {
		return err
	}
	for _, container := range containers {
		if container.ID == id {
			return b.restartContainer(ctx, container.ID)
		}
	}
	return fmt.Errorf("container is not part of the configured Compose project")
}

func (b *Bot) restartContainer(ctx context.Context, id string) error {
	client := b.dockerClient()
	path := "/containers/" + url.PathEscape(id) + "/restart?t=10"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://docker"+path, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
	response.Body.Close()
	if response.StatusCode/100 != 2 {
		return fmt.Errorf("Docker API %s: %s", response.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func (b *Bot) dockerClient() *http.Client {
	client := &http.Client{Timeout: 30 * time.Second, Transport: &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "unix", b.config.DockerSocket)
	}}}
	return client
}

func (b *Bot) listProjectContainers(ctx context.Context, client *http.Client) ([]dockerContainer, error) {
	endpoint := "http://docker/containers/json?all=1"
	filters, _ := json.Marshal(map[string][]string{"label": {"com.docker.compose.project=" + b.config.ComposeProject}})
	endpoint += "&filters=" + url.QueryEscape(string(filters))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		response.Body.Close()
		return nil, fmt.Errorf("Docker API %s: %s", response.Status, strings.TrimSpace(string(body)))
	}
	var containers []dockerContainer
	err = json.NewDecoder(response.Body).Decode(&containers)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(containers, func(i, j int) bool {
		a, c := strings.TrimPrefix(containers[i].Labels["com.docker.compose.service"], ""), strings.TrimPrefix(containers[j].Labels["com.docker.compose.service"], "")
		if (a == "bot") != (c == "bot") {
			return a != "bot"
		}
		return a < c
	})
	return containers, nil
}

func (b *Bot) handleAdminContainers(w http.ResponseWriter, r *http.Request, userID int64) {
	if b.config.AdminTelegramID == 0 || userID != b.config.AdminTelegramID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	switch r.Method {
	case http.MethodGet:
		containers, err := b.listProjectContainers(r.Context(), b.dockerClient())
		if err != nil {
			log.Printf("inspect Docker containers for admin panel: %v", err)
			http.Error(w, "Не удалось получить список контейнеров из Docker: "+err.Error(), http.StatusBadGateway)
			return
		}
		type item struct {
			ID      string `json:"id"`
			Service string `json:"service"`
			Name    string `json:"name"`
			State   string `json:"state"`
		}
		items := make([]item, 0, len(containers))
		for _, c := range containers {
			name := strings.TrimPrefix(strings.Join(c.Names, ", "), "/")
			items = append(items, item{ID: c.ID, Service: c.Labels["com.docker.compose.service"], Name: name, State: c.State})
		}
		writeJSON(w, map[string]any{"project": b.config.ComposeProject, "containers": items})
	case http.MethodPost:
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/containers/")
		parts := strings.Split(path, "/")
		action := ""
		containerID := ""
		if len(parts) == 1 && (parts[0] == "restart" || parts[0] == "delete") {
			action = parts[0]
		} else if len(parts) == 2 && parts[0] != "" && parts[1] == "restart" {
			action, containerID = "restart-one", parts[0]
		} else {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			count, err := 0, error(nil)
			if action == "restart-one" {
				err = b.restartProjectContainer(ctx, containerID)
				if err == nil {
					count = 1
				}
			} else {
				count, err = b.manageProjectContainers(ctx, action)
			}
			if err != nil {
				log.Printf("admin container %s failed after %d: %v", action, count, err)
				return
			}
			log.Printf("admin container %s completed: %d containers", action, count)
		}()
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (b *Bot) handleAdminAccess(w http.ResponseWriter, r *http.Request, userID int64) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]bool{"is_admin": b.config.AdminTelegramID != 0 && userID == b.config.AdminTelegramID})
}

func (b *Bot) sendMessageWithKeyboard(ctx context.Context, chatID int64, message string, keyboard adminMarkup) error {
	payload := struct {
		ChatID      int64       `json:"chat_id"`
		Text        string      `json:"text"`
		ReplyMarkup adminMarkup `json:"reply_markup"`
	}{chatID, message, keyboard}
	return b.telegramJSON(ctx, "sendMessage", payload)
}

func (b *Bot) answerCallback(ctx context.Context, id, message string) error {
	return b.telegramJSON(ctx, "answerCallbackQuery", struct {
		ID   string `json:"callback_query_id"`
		Text string `json:"text,omitempty"`
	}{id, message})
}

func (b *Bot) telegramJSON(ctx context.Context, method string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.config.TelegramBotToken, method)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := b.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return fmt.Errorf("telegram %s: %s: %s", method, response.Status, strings.TrimSpace(string(body)))
	}
	return nil
}
