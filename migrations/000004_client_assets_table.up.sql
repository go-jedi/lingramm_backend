CREATE TABLE IF NOT EXISTS client_assets(
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name_file TEXT NOT NULL, -- Имя файла.
    server_path_file TEXT NOT NULL, -- Путь до файла для сервера.
    client_path_file TEXT NOT NULL, -- Путь до файла для клиента.
    extension VARCHAR(255) NOT NULL, -- Расширение файла.
    quality INTEGER NOT NULL, -- Выставленное качество.
    old_name_file TEXT NOT NULL, -- Имя файла до конвертирования.
    old_extension VARCHAR(255) NOT NULL, -- Расширение файла до конвертирования.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() -- Дата обновления записи.
);