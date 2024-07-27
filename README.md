# supremacy-day-change
discord go bot to alert users about a game's day change

# env variables
Env variables are accessed from the .env file in the root directory.  
The current variables used are:  
BOT_TOKEN - the token that authenticates your discord bot  
ROLE - the role that the bot will send the alert message to  


# commands
/alertHelp - lists the commands and their usage

/setAlert  - ```/setAlert:frequency:startTimestamp:message``` ,   

frequency - the interval at which the alert will fire off,   
startTimestamp: the epoch/unix timestamp to start the alert,   
message: the alert message  
e.g. ```/setAlert:8h:1720571910:hello```  

/stopAlert - stops all set alerts
