# datadogunifi

UniFi Poller Output Plugin for DataDog

## Configuration

```yaml
datadog:
  # How often to poll UniFi and report to Datadog.
  interval: "2m"

  # To disable this output plugin
  disable: false

  # Datadog Custom Options

  # address to talk to the datadog agent, by default this uses the local statsd UDP interface
  # address: "..."

  # namespace to prepend to all data
  # namespace: ""

  # tags to append to all data
  # tags:
  #  - foo
  
  # max_bytes_per_payload is the maximum number of bytes a single payload will contain.
  # The magic value 0 will set the option to the optimal size for the transport
  # protocol used when creating the client: 1432 for UDP and 8192 for UDS.
  # max_bytes_per_payload: 0
  
  # max_messages_per_payload is the maximum number of metrics, events and/or service checks a single payload will contain.
  # This option can be set to `1` to create an unbuffered client.
  # max_messages_per_payload: 0
  
  # BufferPoolSize is the size of the pool of buffers in number of buffers.
  # The magic value 0 will set the option to the optimal size for the transport
  # protocol used when creating the client: 2048 for UDP and 512 for UDS.
  # buffer_pool_size: 0

  # buffer_flush_interval is the interval after which the current buffer will get flushed.
  # buffer_flush_interval: 0
  
  # buffer_shard_count is the number of buffer "shards" that will be used.
  # Those shards allows the use of multiple buffers at the same time to reduce
  # lock contention.
  # buffer_shard_count: 0
  
  # sender_queue_size is the size of the sender queue in number of buffers.
  # The magic value 0 will set the option to the optimal size for the transport
  # protocol used when creating the client: 2048 for UDP and 512 for UDS.
  # sender_queue_size: 0
  
  # write_timeout_uds is the timeout after which a UDS packet is dropped.
  # write_timeout_uds: 5000
  
  # receive_mode determines the behavior of the client when receiving to many
  # metrics. The client will either drop the metrics if its buffers are
  # full (ChannelMode mode) or block the caller until the metric can be
  # handled (MutexMode mode). By default the client will MutexMode. This
  # option should be set to ChannelMode only when use under very high
  # load.
  # 
  # MutexMode uses a mutex internally which is much faster than
  # channel but causes some lock contention when used with a high number
  # of threads. Mutex are sharded based on the metrics name which
  # limit mutex contention when goroutines send different metrics.
  # 
  # ChannelMode: uses channel (of ChannelModeBufferSize size) to send
  # metrics and drop metrics if the channel is full. Sending metrics in
  # this mode is slower that MutexMode (because of the channel), but
  # will not block the application. This mode is made for application
  # using many goroutines, sending the same metrics at a very high
  # volume. The goal is to not slow down the application at the cost of
  # dropping metrics and having a lower max throughput.
  # receive_mode: 0
  
  # channel_mode_buffer_size is the size of the channel holding incoming metrics
  # channel_mode_buffer_size: 0
  
  # aggregation_flush_interval is the interval for the aggregator to flush metrics
  # aggregation_flush_interval: 0
```