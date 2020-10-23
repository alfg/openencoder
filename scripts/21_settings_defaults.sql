INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (1, 'digitaloceanspaces', 9, false);
INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (2, 'sfo2', 8, false);
INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (3, 'outbound', 6, false);
INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (4, 'inbound', 5, false);
INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (5, 'sfo2', 7, false);
INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (6, 'disabled', 10, false);
INSERT INTO public.settings (id, value, settings_option_id, encrypted) VALUES (7, 's3', 14, false);

SELECT setval('settings_id_seq', max(id)) FROM settings;