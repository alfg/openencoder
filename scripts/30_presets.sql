INSERT INTO public.presets (id, name, description, data, active, output) VALUES (5, 'webm_vp9_720p_3000', 'webm_vp9_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libvpx-vp9","-level:v 4.0","-b:v 3000k","-pix_fmt yuv420p","-f webm","-y"]}', true, 'webm_vp9_3000_720.webm');
INSERT INTO public.presets (id, name, description, data, active, output) VALUES (4, 'h264_main_1080p_6000', 'h264_main_1080p_6000', '{"raw":["-vf scale=-2:1080","-c:v libx264","-profile:v main","-level:v 4.2","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 6000k","-maxrate 6000k","-bufsize 6000k","-b:v 6000k","-y"]}', true, 'h264_main_1080p_6000.mp4');
INSERT INTO public.presets (id, name, description, data, active, output) VALUES (3, 'h264_main_720p_3000', 'h264_main_720p_3000', '{"raw":["-vf scale=-2:720","-c:v libx264","-profile:v main","-level:v 4.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 3000k","-maxrate 3000k","-bufsize 3000k","-b:v 3000k","-y"]}', true, 'h264_main_720p_3000.mp4');
INSERT INTO public.presets (id, name, description, data, active, output) VALUES (2, 'h264_main_480p_1000', 'h264_main_480p_1000', '{"raw":["-vf scale=-2:480","-c:v libx264","-profile:v main","-level:v 3.1","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 1000k","-maxrate 1000k","-bufsize 1000k","-b:v 1000k","-y"]}', true, 'h264_main_480p_1000.mp4');
INSERT INTO public.presets (id, name, description, data, active, output) VALUES (1, 'h264_baseline_360p_600', 'h264_baseline_360p_600', '{"raw":["-vf scale=-2:360","-c:v libx264","-profile:v baseline","-level:v 3.0","-x264opts scenecut=0:open_gop=0:min-keyint=72:keyint=72","-minrate 600k","-maxrate 600k","-bufsize 600k","-b:v 600k","-y"]}', true, 'h264_baseline_360p_600.mp4');

SELECT setval('presets_id_seq', max(id)) FROM presets;