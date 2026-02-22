
# Endpoints (nginx conf)
* Go server $SPEEDBALL_DOMEN/speedball_webhook:8443 -> localhsot:8800
* Go server $SPEEDBALL_DOMEN/speedball_notify:8443 -> localhsot:8800
* Go server $SPEEDBALL_DOMEN/speedball-api/v1:8443 -> localhsot:8801


# Dependencies for api and crypto daemon
* Redis

# env:

## Speedball tg-bot daemon
* SPEEDBALL_TG_BOT_TOKEN
* SPEEDBALL_DOMEN

## Speedball API
* REDIS_PASS

