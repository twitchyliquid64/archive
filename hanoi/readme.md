Hanoi
=========

Your infrastructure watchdog.

### Overview

Hanoi gathers its data by running `sensors` according to their configuration.

Sensor data is then interpreted by `rules`, which tell Hanoi when something is bad or not.

If something is bad, `pager rules` fire.

For organisation sake, sensors, pager rules, and rules are grouped. Typically, they are grouped by the system you are monitoring.

### Pushtart deployment supported.

Hanoi was built to run on pushtart.

### Launching

Run on the command line like: `python main.py <server_port> <path_to_config_file>`

### Configuration file format

```json
{
  "groups": {
    "<group name>": {
      "sensors": [
        ...
      ],
      "rules": [
        ...
      ],
      "paging_rules": [
        ...
      ]
    },
  },
  "gmail": {
    "username": "<your username>@gmail.com",
    "password": "<generated app password>"
  }
}
```

Sensor item format:

The only required fields are `type` and `name`. Name is used by rules to source information from the sensor.
The type corresponds to a sensor script in `sensors/`. Other fields in this structure are dependent on the type of sensor.

The field `explanation_exp` is optional, but will be evaluated (as JS) and displayed in the UI.

```json
{
  "type": "latency",
  "dest": "google.com",
  "interval": 300,
  "name": "google_latency",
  "explanation_exp": "String(sensor.last_result.val) + \" :)\"",
}
```

Rule item format:

Notice the name match between the sensor item, and the name of the above sensor.
As you can guess, `condition` is python code. If it evaluates to true, the rule is passing.
In this case, the rule checks the sensor is functioning, and the latency is less than 150ms.

The field `explanation_exp` is optional, but will be evaluated (as JS) and displayed in the UI.

```json
{
  "name": "Network Latency",
  "type": "latest",
  "sensors": {
    "s": "google_latency"
  },
  "condition": "s.ok and s.last_result['val'] < 150",
  "explanation_exp": "(rule.ok && !rule.noop) ? \"OK - Latency less than 150ms\" : (rule.noop ? \"Sensor not ready\" : \"Latency exceeds 150ms\")"
}
```

Paging rule format:

Fairly simple. Items in the `rules` array should correspond to the names of rules. If any rules in the list fail, the pager is triggered.
The pager will not fire again unless the rules all succeed again before failing again. `hysteresis` is an optional field which specifies
a minimum number of seconds before paging again - to avoid waking up to hundreds of emails if the rule is succeeding then failing etc.


```json
{
  "name": "Latency too high",
  "rules": [
    "Network Latency"
  ],
  "hysteresis": 80,
  "page": {
    "type": "GMAIL",
    "address": "my_email_@gmail.com"
  }
}
```

At the moment, only sending emails from a gmail account is supported. You will need to put the following in the root of the configuration file:

```json
"gmail": {
  "username": "<your username>@gmail.com",
  "password": "<generated app password>"
}
```

### Almost finished.

 - [ ] Implement statistical (rate of change, histeresis) rule types
 - [x] Display paging rules in UI
 - [ ] Make homepage flashing status coupled to triggering a pager
 - [x] Validate paging_rules
 - [x] Implement paging system + config.
