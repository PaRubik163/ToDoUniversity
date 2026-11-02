-- Invalidate pending codes issued by the older, lower-entropy generator.
DELETE FROM auth_link_codes;
