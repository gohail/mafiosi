#### **GAME SERVER ACTIONS**

**Connect to ->** `ws://server_adress:8080/`

successes response:

````
{
"view" : "START_VIEW",
"err"  : "",
"data" : {}
}
````
---
**START_VIEW** actions:

- *_CREATE_GAME_*

expected client json:
````
{
"action" : "CREATE_GAME"
}
````

successes server response:
````
{
"view" : "REQ_NAME",
"err"  : "",
"data" : null
}
````

expected client json:
````
{
"name" : "user name" (not empty)
}
````
OR:

````
{
"action" : "CANCEL"
}
````
(CANCEL will return you to "START_VIEW" loop)
