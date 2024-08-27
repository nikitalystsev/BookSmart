insert into bs.reader
values ('75919792-c2d9-4685-92b2-e2a80b2ed5be', 'Randall C. Jernigan', '79314562376', 25,
        '$2a$10$8APnhcfxoGxXGdNHSdBEaebuwcIkjwEnSHOIv.xu9bmROkpCRLTJS', 'Reader'), -- пароль: sdgdgsgsgd
       ('5818061a-662d-45bb-a67c-0d2873038e65', 'Jesse M. Flores', '72443564633', 20,
        '$2a$10$2cYeMgl8fjH76HjIm54enOuHUiV3qzV81jdVJLLNCQbo2zXc9jija', 'Reader'), -- пароль: qwresdfdsf
       ('3885b2d3-ef6e-4f62-8f86-d1454d108207', 'Mitrofan Bogdanov', '76867456521', 40,
        '$2a$10$KxEprnJtxnL./4Zts.IP3uOQfGktXZXTp1BmvMKZxyDSJoIm4hmt6', 'Reader'), -- пароль: hghhfnnbdd
       ('6800b3ee-9810-450e-9ca5-776aa1c6191d', 'Peter Zuev', '32534523451', 13,
        '$2a$10$GjKIYnr6wRohYWkUhmlPhO5uza1zvudS9rWeydAv1yzEW0GfTOAme', 'Reader'), -- пароль: rtjhhhgffr
       ('8d9b001f-5760-4c40-bc60-988e0ca54d18', 'Vasilisa Agapova', '73453562423', 36,
        '$2a$10$sQZzp5BlhAvTMc/AIzAUS.PVuAxxH/rVmNfv.W73RhdxH7xSdbyQy', 'Reader'), -- пароль: gfjkjdgffy
       ('362b79f6-d671-404a-b1a0-5a655aebc1b6', 'Лысцев Никита Дмитриевич', '89314022581', 21,
        '$2a$10$xDzRFS0ClhEcosyFVQEPCev8AXakZyYau4Hk8iN3dyTXJYXUj1coO', 'Admin');

insert into bs.lib_card
values ('e71af5a9-dd02-4f00-982e-ec58908ec5bd', '75919792-c2d9-4685-92b2-e2a80b2ed5be', '4654645456328', 365,
        '2024-07-26', true),
       ('5411019a-7dbf-4dbb-b621-784176da6ec5', '5818061a-662d-45bb-a67c-0d2873038e65', '3455787683242', 365,
        '2024-05-13', true),
       ('0614f5de-ac79-4bab-995f-99944a7ca4db', '3885b2d3-ef6e-4f62-8f86-d1454d108207', '7945544456734', 365,
        '2022-07-26', false),
       ('894f6d5c-f81a-46c0-98aa-d7a90aafd93e', '6800b3ee-9810-450e-9ca5-776aa1c6191d', '5435645425466', 365,
        '2023-03-05', false);

insert into bs.reservation
values ('89ff79cd-5ef9-4553-9dac-b3fc2954048c', '75919792-c2d9-4685-92b2-e2a80b2ed5be',
        'f01107fb-4f7a-4f37-ba1e-6c6012c5203c', '2024-08-26', '2024-09-09', 'Issued'),
       ('97f918b0-b8f6-4bd4-b0e0-9d25971b9af2', '75919792-c2d9-4685-92b2-e2a80b2ed5be',
        '43f45552-4a95-4f12-864b-e1d8bfa30b8d', '2024-08-26', '2024-09-09', 'Issued'),
       ('bdf16862-6289-4f21-aff7-2e5fdb9edfed', '5818061a-662d-45bb-a67c-0d2873038e65',
        'b33b30c8-254e-45f2-8314-0b93a6b8c561', '2024-05-14', '2024-05-28', 'Expired');

create role administrator;

grant usage on schema bs to administrator;
grant all privileges on all tables in schema bs to administrator;

create user admin_user with password 'admin';

grant administrator to admin_user;