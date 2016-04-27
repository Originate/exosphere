So let's build our "todo-search-service"
using an already existing off-the-shelf search service.

```
cd ~/todo-app
exo link service
```

We enter into the wizard:

<table>
  <tr>
    <th>name of the service in our app</th>
    <td>todo-search-service</td>
  </tr>
  <tr>
    <th>URL of the service</th>
    <td>https://github.com/originate/fulltext-search-service</td>
  </tr>
  <tr>
    <th>version</th>
    <td>1.2</td>
  </tr>
</table>

Exosphere now adds the new service to the
[application configuration](03_app_config.md):

__~/todo-app/application.yml__

```yml
name: Todo application
...

services:
  - web:
    - location: web-server
  - todo:
    - location: todo-service
  - search:
    - location: https://github.com/originate/fulltext-search-service
    - version: 1.2
```



