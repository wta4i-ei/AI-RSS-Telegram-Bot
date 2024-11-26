# AI-RSS-Telegram-Bot

Bot for Telegram that gets, summarizes and posts news to a channel.

## Features

- Fetching articles from RSS feeds
- Article summaries powered by GPT
- Admin commands for managing sources

## Demonstration
<details>
<summary>Click to view the demonstration</summary>
Here’s how the bot works in action:

### 1. Posting posts to the channel

The bot publishes the title of the article, a summary and a link to the original article

![Posting posts](./images/img_2.png)

### 2. Adding a source

![/addsource](./images/img_1.png)

### 3. Deleting a source

![/deletesource](./images/img_3.png)

### 4. Setting the priority of the source

![/setpriority](./images/img_4.png)

### 5. Displaying a list of sources


![/listsources](./images/img.png)

</details>

## Technology Stack

- **Programming Language:** Go
- **Database:** PostgreSQL
- **AI Summarization:** OpenAI GPT-3.5 Turbo
- **Configuration:** HCL and environment variables
- **Containerization:** Docker and Docker Compose


## **How to Run the Bot**

### **1. Clone the Repository**

```bash
git clone https://github.com/wta4i-ei/AI-RSS-Telegram-Bot
cd project
```

---

### **2. Configure the Bot**

You can configure the bot using either **environment variables** (`.env` file) or **HCL configuration files**.

#### **Environment Variables**

Create a `.env` file in the project root with the following variables:

```env
PROXY_URL=your_proxy_url
NFB_TELEGRAM_BOT_TOKEN=your_telegram_bot_token
NFB_TELEGRAM_CHANNEL_ID=your_channel_id
NFB_DATABASE_DSN=postgres://user:password@host:port/dbname?sslmode=disable
NFB_FETCH_INTERVAL=10m
NFB_NOTIFICATION_INTERVAL=1m
NFB_OPENAI_KEY=your_openai_api_key
NFB_OPENAI_PROMPT="Summarize the article concisely in 3 sentences"
```

#### **HCL Variables**

Alternatively, create an HCL config file in one of the following locations:

- `./config.hcl`
- `./config.local.hcl`
- `$HOME/.config/AI-RSS-Telegram-Bot/config.hcl`

Example `config.hcl`:

```hcl
proxy_url = "your_proxy_url"
telegram_bot_token = "your_telegram_bot_token"
telegram_channel_id = "your_channel_id"
database_dsn = "postgres://user:password@host:port/dbname?sslmode=disable"
fetch_interval = "10m"
notification_interval = "1m"
openai_key = "your_openai_api_key"
openai_prompt = "Summarize the article concisely in 3 sentences"
```

---

### **3. Start the Bot Using Docker Compose**

Run the following command to build and start the bot:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

---

### **4. Verify Docker Containers**

Check that the containers are running:

```bash
docker ps
```

---

### **5. Install Migration Tool**

Install `goose`, a database migration tool:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

### **6. Run Migrations**

Navigate to the migrations directory and apply migrations to the database:

```bash
goose postgres "host=localhost user=postgres database=news_feed_bot password=postgres sslmode=disable" up
```

---

### **7. Start the Bot**

Navigate to the `cmd` directory and run the bot:

```bash
cd cmd
go run main.go
```

---

### **8. Manage Sources Using Bot Commands**

Use the following commands in your Telegram bot to manage RSS sources:

- **Add a Source**  
  `/addsource 'name' 'url' 'priority'`  
  Example: `/addsource "Tech News" "https://example.com/rss" "1"`

- **Set Priority for a Source**  
  `/setpriority 'id' 'priority'`  
  Example: `/setpriority 1 10`

- **Delete a Source**  
  `/deletesource 'id'`  
  Example: `/deletesource 1`

- **Get Information About a Source**  
  `/getsource 'id'`  
  Example: `/getsource 1`

- **List All Sources**  
  `/listsources`

---

## **Done!**

The bot is now configured and running. It will fetch news articles based on the `fetch_interval` and post them to your Telegram channel at the `notification_interval`. 🎉

