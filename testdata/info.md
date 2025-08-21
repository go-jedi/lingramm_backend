### generate swagger:
- `swag init -g cmd/app/main.go`

### migrations:

#### create:
- `migrate create -ext sql -dir migrations -seq set_timezone`
- `migrate create -ext sql -dir migrations -seq users_table`
- `migrate create -ext sql -dir migrations -seq user_create_function`
- `migrate create -ext sql -dir migrations -seq levels_table`
- `migrate create -ext sql -dir migrations -seq client_assets_table`
- `migrate create -ext sql -dir migrations -seq user_balances_table`
- `migrate create -ext sql -dir migrations -seq user_balances_index`
- `migrate create -ext sql -dir migrations -seq balance_transaction_events_table`
- `migrate create -ext sql -dir migrations -seq balance_transactions_table`
- `migrate create -ext sql -dir migrations -seq balance_transactions_index`
- `migrate create -ext sql -dir migrations -seq currency_rates_table`
- `migrate create -ext sql -dir migrations -seq users_blacklist_table`
- `migrate create -ext sql -dir migrations -seq admins_table`
- `migrate create -ext sql -dir migrations -seq achievement_assets_table`
- `migrate create -ext sql -dir migrations -seq award_assets_table`
- `migrate create -ext sql -dir migrations -seq achievements_table`
- `migrate create -ext sql -dir migrations -seq achievement_conditions_table`
- `migrate create -ext sql -dir migrations -seq user_achievements_table`
- `migrate create -ext sql -dir migrations -seq user_stats_table`
- `migrate create -ext sql -dir migrations -seq text_contents_table`
- `migrate create -ext sql -dir migrations -seq text_translations_table`
- `migrate create -ext sql -dir migrations -seq daily_tasks_table`
- `migrate create -ext sql -dir migrations -seq user_daily_tasks_table`
- `migrate create -ext sql -dir migrations -seq notifications_type`
- `migrate create -ext sql -dir migrations -seq notifications_table`
- `migrate create -ext sql -dir migrations -seq notifications_index`
- `migrate create -ext sql -dir migrations -seq subscriptions_table`
- `migrate create -ext sql -dir migrations -seq subscription_history_table`
- `migrate create -ext sql -dir migrations -seq subscription_create_function`
- `migrate create -ext sql -dir migrations -seq subscriptions_exists_function`
- `migrate create -ext sql -dir migrations -seq xp_events_table`
- `migrate create -ext sql -dir migrations -seq xp_events_index`
- `migrate create -ext sql -dir migrations -seq xp_event_create_function`
- `migrate create -ext sql -dir migrations -seq sync_user_stats_from_xp_events_function`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_table`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_index`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_worker_state_table`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_worker_state_index`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_applied_events_table`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_process_batch_function`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_top_week_get_function`
- `migrate create -ext sql -dir migrations -seq leaderboard_weeks_top_week_for_user_get_function`
- `migrate create -ext sql -dir migrations -seq user_level_history_table`
- `migrate create -ext sql -dir migrations -seq user_level_history_index`
- `migrate create -ext sql -dir migrations -seq back_fill_missing_level_history_function`
- `migrate create -ext sql -dir migrations -seq get_level_info_function`
- `migrate create -ext sql -dir migrations -seq languages_table`

#### execute:
- `migrate -database postgresql://admin:test@localhost:54320/lingvogramm_db?sslmode=disable -path migrations up`
- `migrate -database postgresql://admin:test@localhost:54320/lingvogramm_db?sslmode=disable -path migrations down`

#### build application:
- `go build -ldflags="-s -w" -trimpath -buildvcs=false -o app cmd/app/main.go`

#### run application in systemd:
- `cd /etc/systemd/system`
- `создать lingvogramm_backend.service`
- `sudo systemctl daemon-reload`
- `sudo systemctl start lingvogramm_backend.service`
- `sudo systemctl status lingvogramm_backend.service`
- `sudo systemctl enable lingvogramm_backend.service`

#### Включить порт:
- `sudo ufw allow 50051/tcp`
- `sudo ufw reload`
- если при выполнении команды: sudo ss -tuln | grep 50051 у вас показывается:
  tcp    LISTEN  0       4096         127.0.0.1:50051        0.0.0.0:*
  , то это указывает, что сервис будет доступен только внутри сервера через localhost.
- Если нужно, чтобы можно было отправлять запросы из внешних источников:
- в сервисе указываем при запуске http сервера :50051
- выполняем sudo ss -tuln | grep 50051 и должно быть в ответе такой результат:
  tcp    LISTEN  0       4096           0.0.0.0:50051        0.0.0.0:*

#### docker build and push to docker hub:
- `docker build -t gojedi/lingvogramm_backend:latest .`
- `docker push gojedi/lingvogramm_backend:latest`

#### remove all local branch without main:
- `git branch | grep -v "main" | xargs git branch -D`

#### cancel commit:
- `git reset --soft HEAD~1`


## Про cookie в fiber v3:
Вот подробное объяснение каждого поля структуры `Cookie` в **Golang Fiber v3**:

---

### 🔐 Основные поля:

| Поле              | Тип         | Описание                                                                                          |
| ----------------- | ----------- | ------------------------------------------------------------------------------------------------- |
| **`Name`**        | `string`    | Название куки. Например, `"session_id"` или `"user_token"`.                                       |
| **`Value`**       | `string`    | Значение куки. Часто — зашифрованный ID сессии или токен.                                         |
| **`Expires`**     | `time.Time` | Время окончания действия куки. После этой даты браузер удалит её.                                 |
| **`MaxAge`**      | `int`       | Время жизни куки в **секундах**. Альтернатива `Expires`. Если `MaxAge <= 0`, кука не сохраняется. |
| **`SessionOnly`** | `bool`      | Указывает, должна ли кука жить только в рамках сессии (закрытие вкладки = удаление).              |

---

### 📍 Область действия:

| Поле         | Тип      | Описание                                                                                                                               |
| ------------ | -------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| **`Path`**   | `string` | Путь, на который кука распространяется. Например, `/admin` означает, что кука будет доступна только для URL, начинающихся на `/admin`. |
| **`Domain`** | `string` | Домен, для которого доступна кука. Например, `.example.com` позволит использовать куку на всех поддоменах.                             |

---

### 🛡️ Безопасность:

| Поле           | Тип      | Описание                                                                                                                                                                                                                               |
| -------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`Secure`**   | `bool`   | Если `true`, кука будет передаваться **только по HTTPS**.                                                                                                                                                                              |
| **`HTTPOnly`** | `bool`   | Если `true`, кука **не доступна из JavaScript** (например, через `document.cookie`). Защита от XSS.                                                                                                                                    |
| **`SameSite`** | `string` | Защита от CSRF. Возможные значения: <br> 🔹 `"Strict"` — кука передаётся только при прямом заходе;<br> 🔹 `"Lax"` — кука передаётся при переходах по ссылкам;<br> 🔹 `"None"` — разрешает кросс-доменные куки (требует `Secure=true`). |

---

### 🧪 Продвинутые поля:

| Поле              | Тип    | Описание                                                                                                                                                                                             |
| ----------------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`Partitioned`** | `bool` | Разделённое хранилище куки для разных топ-уровневых сайтов. Введено в Chrome как эксперимент — полезно в `iframe` и при изоляции third-party cookies. <br>Необходим `SameSite=None` и `Secure=true`. |

---

### 📝 Пример установки полной куки:

```go
c.Cookie(&fiber.Cookie{
    Name:        "token",
    Value:       "abc123",
    Expires:     time.Now().Add(24 * time.Hour),
    Path:        "/",
    Domain:      "example.com",
    MaxAge:      86400, // 24 часа
    Secure:      true,
    HTTPOnly:    true,
    SameSite:    "Lax",
    SessionOnly: false,
})
```

Если тебе нужно на практике сравнение `Expires` и `MaxAge`, или отличие `SessionOnly` и `Expires`, — могу объяснить через примеры.
