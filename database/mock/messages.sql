insert into chat (id, name) values ('c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'neastwood0');
insert into chat (id, name) values ('fd703f8a-6674-475f-a970-eea68504d815', 'cambrosi1');
insert into chat (id, name) values ('aa504f61-e421-4756-813c-17d2ee422e02', 'cchaffer2');

insert into chat_user (chat_id, user_id) values ('c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', '751038f9-7781-4294-b9dc-19841c03c044');
insert into chat_user (chat_id, user_id) values ('c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d');
insert into chat_user (chat_id, user_id) values ('fd703f8a-6674-475f-a970-eea68504d815', '751038f9-7781-4294-b9dc-19841c03c044');
insert into chat_user (chat_id, user_id) values ('fd703f8a-6674-475f-a970-eea68504d815', 'fd703f8a-6674-475f-a970-eea68504d815');
insert into chat_user (chat_id, user_id) values ('aa504f61-e421-4756-813c-17d2ee422e02', '751038f9-7781-4294-b9dc-19841c03c044');
insert into chat_user (chat_id, user_id) values ('aa504f61-e421-4756-813c-17d2ee422e02', 'aa504f61-e421-4756-813c-17d2ee422e02');

insert into message (id, sender_id, chat_id, message, created_at) values ('c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', '751038f9-7781-4294-b9dc-19841c03c044', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'ciao!', '2020-01-01 00:00:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('fd703f8a-6674-475f-a970-eea68504d815', '751038f9-7781-4294-b9dc-19841c03c044', 'fd703f8a-6674-475f-a970-eea68504d815', 'ciao!', '2020-01-01 00:00:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('aa504f61-e421-4756-813c-17d2ee422e02', '751038f9-7781-4294-b9dc-19841c03c044', 'aa504f61-e421-4756-813c-17d2ee422e02', 'ciao!', '2020-01-01 00:00:00');

insert into message (id, sender_id, chat_id, message, created_at) values ('0ee56e1c-70ec-449d-a5d7-174b1f3ca9da', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'ciao, come va?', '2020-01-01 00:01:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('0fdbabb6-cb73-4a8a-b81d-74a58eb83a86', 'fd703f8a-6674-475f-a970-eea68504d815', 'fd703f8a-6674-475f-a970-eea68504d815', 'ciao, come va?', '2020-01-01 00:01:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('588a662a-6e14-4e62-a4e2-628feb7e9d25', 'aa504f61-e421-4756-813c-17d2ee422e02', 'aa504f61-e421-4756-813c-17d2ee422e02', 'ciao, come va?', '2020-01-01 00:01:00');

insert into message (id, sender_id, chat_id, message, created_at) values ('0ea130df-8eff-4d91-8d39-3cafba9f43f7', '751038f9-7781-4294-b9dc-19841c03c044', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'bene dai, te?', '2020-01-01 00:02:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('e94e7b58-feba-49c4-bba8-f09715140599', '751038f9-7781-4294-b9dc-19841c03c044', 'fd703f8a-6674-475f-a970-eea68504d815', 'bene dai, te?', '2020-01-01 00:02:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('04e759ac-c6bc-4780-bb23-c8d6563670f4', '751038f9-7781-4294-b9dc-19841c03c044', 'aa504f61-e421-4756-813c-17d2ee422e02', 'bene dai, te?', '2020-01-01 00:02:00');

insert into message (id, sender_id, chat_id, message, created_at) values ('f01364f4-6ce8-45ce-93aa-d92fc2a933d8', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'c0c11cc4-bfab-4a7c-8fe2-43fc477d527d', 'tutto bene.', '2020-01-01 00:03:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('32e75458-a7a8-4047-9b4a-eb0864517f1b', 'fd703f8a-6674-475f-a970-eea68504d815', 'fd703f8a-6674-475f-a970-eea68504d815', 'tutto bene.', '2020-01-01 00:03:00');
insert into message (id, sender_id, chat_id, message, created_at) values ('93275687-c306-4624-9393-486176ea8cb1', 'aa504f61-e421-4756-813c-17d2ee422e02', 'aa504f61-e421-4756-813c-17d2ee422e02', 'tutto bene.', '2020-01-01 00:03:00');


