## Analytics Dashboard

I took a quick look into the available files from stackoverflow and decide to use the following ones:
- Posts
- Users
- Tags
- Votes
- + a few of them that relate to each other...

With these files I think I should be able to add the following data to the dashboard:
- users -> post count
- users -> with the highest votes
- posts -> with the highest votes
- tags -> most seeing tags
- histograms of users x posts made

There's a bunch of data that I can show here but for now these are fine...

### System overview

Dashboard initial screen -> select api [golang, python] -> data screen showing the proposed data

### Frontend stack
Please do not expect any fancy things here, frontend is not my thing
-> Svelte + tailwind

### Backend stack
-> Golang api -> Gin
-> Python api -> FastApi
-> Postgresql DB