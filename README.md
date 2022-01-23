Dropin replacement for the [collectd Redis plugin](https://www.collectd.org/wiki/index.php/Plugin:Redis) using [collectd Exec](https://collectd.org/wiki/index.php/Plugin:Exec).

The [collectd Redis Plugin](https://www.collectd.org/wiki/index.php/Plugin:Redis) submits wrong metrics when collectd runs with multithreading (see https://github.com/collectd/collectd/issues/3941). Multithreading is the default and is highly recommended to prevent plugins from blocking other plugins.


## Configuration

Get the binary and place it on a server running collectd. Make it executable.

This plugin uses [collectd Exec](https://collectd.org/wiki/index.php/Plugin:Exec). So you need an `exec` configuration. Then you can setup several instances to be monitored. The following example shows how a `Redis` Plugin config can be replaced with the `Exec` setup:

### Redis Plugin

```
<LoadPlugin redis>
  Globals false
</LoadPlugin>

<Plugin redis>
  <Node "caching">
    Host "localhost"
    Port "6379"
  </Node>
  <Node "queue">
    Host "localhost"
    Port "6380"
  </Node>
</Plugin>
```

### Exec Plugin

```
<LoadPlugin exec>
  Globals false
</LoadPlugin>

<Plugin exec>
  Exec "collectd:collectd" "/opt/collectd/collectd-redis" "caching:localhost:6379"
  Exec "collectd:collectd" "/opt/collectd/collectd-redis" "queue:localhost:6380"
</Plugin>
```

## Connection string

The plugin takes a connection string in the format `<name>:<host/ip>:<port>[:<password>]`.

The password is optional and the part of the string in square bracktes is not required if the redis instances does not use password authentication.
