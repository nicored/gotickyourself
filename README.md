GO Tick Yourself!!
=================

# How to install

1. Instal golang 1.8
2. Then run the following command to download and build the binary for your system
```shel
    go get github.com/nicored/gotickyourself/gty/...
```
3. Go find the binary and use it.


# What you can do

```shell
    # Configure
    $ gty init
    (Not implemented) $ gty settings # show all settings
    (Not implemented) $ gty settings --reset # resets everything, think twice
    
    # Update roles, projects and tasks list
    $ gty update
    
    # Projects settings
    $ gty projects # List all projects except hidden ones
    $ gty projects -f "McDonalds" # Filter projects by name
    $ gty projects -c "Origin" # Filter projects by client name
    $ gty projects -f "McDonalds" -c "Origin" # Filter projects by name and client name
    
    # Tasks
    $ gty tasks # List all tasks
    $ gty tasks -f development # Filter tasks by name
    $ gty tasks -c "my client" -p "my project" -f "my task" # Filter by client project and task
    $ gty tasks default 1234 default # Default task to ID 1234
    
    # Task Alias
    $ gty tasks alias # List all tasks with an alias
    $ gty tasks alias add dev 12345  # Create "dev" alias for task 12345
    $ gty tasks alias rm 12345 # Remove alias for task 12345
    $ gty tasks alias rm -f # Remove alias for all tasks
    
    # Log new entries
    $ gty log # Automatically log to default task for the remaining hours of the day
    $ gty log -n "I did this today" # Same as log, but with notes
    $ gty log dev # Automatically logs to 'dev' task for the ramaining hours of the day
    $ gty log 2.0 # Logs to default task for today
    $ gty log 2.0 dev # Logs to task associated with alias "dev"
    $ gty log 2.0 -n "Today I did this" # Logs to default task with notes
    $ gty log 2.0 dev -n "Today I developed an app" # Logs to task associated with alias dev, with notes
    $ gty log 3.0 -d yesterday # Log entry for yesterday
    $ gty log 3.0 -d 2017-08-03 -n "I did something" # Log entry for the
    
    # Automatically log for a period of time
    (Not implemented) $ gty log 3 days # Log for the last 3 working days (default task)
    (Not implemented) $ gty log week # Log for the beginning of the week until today (inc) (default task)
    (Not implemented) $ gty log fortnight # Log for 2 weeks (2w prior to the current/previous Monday) (default task)
    (Not implemented) $ gty log 3 weeks # Log for 3 weeks (3w prior to the current/previous Monday) (default task)
    (Not implemented) $ gty log dev week # Log for the week to task associated with alias "dev"
    
    # Automatically round up time
    # If you're a bit short in hours, the round command
    # will round hours in your logs to meet the objective (number of hours per week)
    (Not implemented) $ gty round week
    (Not implemented) $ gty round today
    
    # Total hours
    $ gty sum # for today
    $ gty sum today # for today
    $ gty sum week # for the entire week from Monday
    $ gty sum month # for the entire month
    $ gty sum fortnight # for the 2 previous weeks before the previous monday
    $ gty sum 3 weeks # for the 3 previous weeks before the previous monday
    $ gty sum 3 days # for the last 3 days
    
    # List entries
    $ gty ls # List today's entries
    $ gty ls week # List week's entries
```

# How do you start

First you must initialise

```shell
    $ gty init
    > Username: youemail@domain.com
    > Password: yourPassword
    ...
```

Then, you can select a default task. The default task is useful when you do not want to specify
a task to log your time against.

Let's find a task! The following will list all tasks matching your filter.
Each task will have an ID next to it.

```shell
    $ gty tasks -c "My company" -p "Team 1" -f "development" # I filter my tasks by client, project and task name
```

Once you've got the ID, let's say 1234, you can set the task as default by running the following:

```shell
    $ gty tasks default 1234
```

You can also add aliases to your tasks. For instance, HR in my company added the following tasks:

```shell
    $ gty tasks -f leave

	Administration & Overheads:
		333456 - Annual Leave
		333457 - Sick Leave
```

And I want to add aliases for both Annual and sick leaves:

```shell
    $ gty tasks add alias holidays 333456
    $ gty tasks add alias sick 333457
```

And from now on, whenever I am sick, I can run log sick and it will create a new log entry
to the sick leave task:

```shell
    $ gty log sick
```

I would also recommend to check the ~/.gty/settings.yml file to set the number of hours per week. The default is 40,
and there are 2 non-working days, so it will log 8 hours for the day. 

Something cool too, say you've logged 3 hours already during the day, and you've got to rush outside the office,
you can run 'log' by itself or followed by an alias and it will log 5 hours to the default or aliased task.


# TODO
- Write tests
- gty settings # show all settings
- gty settings --reset # resets everything, think twice
- Automatically log time for a given period
- Round up times
