INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (1, 'AWS_ACCESS_KEY', 'AWS Access Key (Required for S3)', 'AWS Access Key', true);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (2, 'AWS_SECRET_KEY', 'AWS Secret Key (Required for S3)', 'AWS Secret Key', true);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (3, 'DIGITAL_OCEAN_ACCESS_TOKEN', 'Digital Ocean Access Token (Required for Machines)', 'Digital Ocean Access Token', true);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (4, 'SLACK_WEBHOOK', 'Slack Webhook for notifications', 'Slack Webhook', true);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (5, 'S3_INBOUND_BUCKET', 'S3 Inbound Bucket', 'S3 Inbound Bucket', false);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (6, 'S3_INBOUND_BUCKET_REGION', 'S3 Inbound Bucket Region', 'S3 Inbound Bucket Region', false);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (7, 'S3_OUTBOUND_BUCKET', 'S3 Outbound Bucket', 'S3 Outbound Bucket', false);
INSERT INTO public.settings_option (id, name, description, title, secure) VALUES (8, 'S3_OUTBOUND_BUCKET_REGION', 'S3 Outbound Bucket Region', 'S3 Outbound Bucket Region', false);