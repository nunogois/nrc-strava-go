# nrc-strava-go

Sync your Nike Run Club activities to Strava.

## Setup

Create a `.env` file on this folder with this format:

```
START=
NIKE_CLIENT_ID=
NIKE_REFRESH_TOKEN=
STRAVA=
```

### START

START should be the starting date for the search as **YYYY-MM-DD**, like for example **2021-06-01**.

### NRC

1.  Go to [nike.com](https://www.nike.com/) and login;
2.  Open your browser's console by pressing **F12** and open the **Application** tab;
3.  Check **Storage > Local Storage > https://unite.nike.com > com.nike.commerce.nikedotcom.web.credential**;
4.  On the bottom of the screen you will see the corresponding object. You will need to copy **clientId** and **refresh_token** to the corresponding fields on the `.env` file;
5.  **Don't** logout as this will invalidate your token!

### Strava

> WIP

## Run

`go run .`

## TODO

- [ ] Make this work;
- [x] Use refresh tokens whenever possible;
- [ ] Add more info to README;
- [ ] Adapt to use [Viper](https://github.com/spf13/viper);
- [ ] Proper error handling;
- [ ] Refactor requests to an API file;
- [ ] Code review;

## Similar projects

- https://github.com/alexpryshchepa/nrc2strava
- https://github.com/opierre/NRCToStrava
- https://github.com/ygina/nike-strava
- https://github.com/yihong0618/running_page

## Other notes

- https://developers.strava.com/docs/getting-started/#account
