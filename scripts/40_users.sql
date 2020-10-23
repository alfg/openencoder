INSERT INTO public.users (id, username, password, role, force_password_reset, active) VALUES (1, 'admin', '$2a$04$FxaVhOgeUazmjfhe4eGrgeFx/Dm3nyw0/so4k.pPSsVDj.7lZmJDW', 'admin', true, true);

SELECT setval('users_id_seq', max(id)) FROM users;