# nrc-strava-go

Sync your Nike Run Club activities to Strava.

## Setup

Create a `.env` file on this folder with this format:

```
START=

NIKE_CLIENT_ID=
NIKE_REFRESH_TOKEN=

STRAVA_CLIENT_ID=
STRAVA_CLIENT_SECRET=
STRAVA_REFRESH_TOKEN=
```

### START

**START** should be the starting date for the search as **YYYY-MM-DD**, like for example **2021-06-01**.

### NRC

1.  Go to [nike.com](https://www.nike.com/) and login;
2.  Open your browser's console by pressing **F12** and open the **Application** tab;
3.  Check **Storage > Local Storage > https://unite.nike.com > com.nike.commerce.nikedotcom.web.credential**;
4.  On the bottom of the screen you will see the corresponding object. You will need to copy **clientId** and **refresh_token** to the corresponding fields on the `.env` file;
5.  **Don't** logout as this will invalidate your token!

### Strava

1. Go to [Strava API Settings](https://www.strava.com/settings/api);
2. Create a new app. I suggest the following as an example:

**Application Name:** NRC Go

**Website:** https://github.com/nunogois/nrc-strava-go

**Authorization Callback Domain:** nrc-strava-go.nunogois.com

When it asks for the icon, feel free to use the **icon.png** included in this repo.
These are only recommendations, you can use something entirely different and it should still work.

3. Copy **Client ID** and **Client Secret** to the corresponding fields on the `.env` file;
4. Go to https://www.strava.com/oauth/authorize?client_id=CLIENT_ID&response_type=code&redirect_uri=https://nrc-strava-go.nunogois.com&approval_prompt=force&scope=activity:write where you will need to replace **CLIENT_ID** with your own and **nrc-strava-go.nunogois.com** with your callback domain, if you set something different.
5. Authorize the app and you might be redirect to a non-functionable website. That's OK though, from this new URL we need to grab the **code** value, before `&scope=read,activity:write`.
6. We will need to make a request. For this you can use cURL or something like Postman, Insomnia, etc:

```sh
curl -X POST https://www.strava.com/api/v3/oauth/token \
  -d client_id=CLIENT_ID \
  -d client_secret=CLIENT_SECRET \
  -d code=CODE \
  -d grant_type=authorization_code
```

If you're using Postman or Insomnia, for example, you can simply create a new **POST** request to https://www.strava.com/api/v3/oauth/token. The values can be sent on the body as JSON:

```JavaScript
{
	"client_id": "CLIENT_ID",
	"client_secret": "CLIENT_SECRET",
	"code": "CODE",
	"grant_type": "authorization_code"
}
```

Make sure to replace **CLIENT_ID**, **CLIENT_SECRET** and **CODE** accordingly. You should have all of this information from the previous steps.

From the response, you will need to grab the **refresh_token** so you can copy it to **STRAVA_REFRESH_TOKEN** on the `.env` file.

## Run

`go run .`

## TODO

- [x] Make this work;
- [x] Use refresh tokens whenever possible;
- [x] Add more info to README;
- [ ] Adapt to use [Viper](https://github.com/spf13/viper);
- [ ] Refactor requests to an API file;
- [ ] Proper error handling;
- [ ] Code review;
- [ ] Add heartbeat, elevation and other relevant info;

## Similar projects

- https://github.com/alexpryshchepa/nrc2strava
- https://github.com/opierre/NRCToStrava
- https://github.com/ygina/nike-strava
- https://github.com/yihong0618/running_page

## Other notes

- https://developers.strava.com/docs/getting-started/#account
