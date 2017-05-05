SystemdConsumer
===============

NOTICE: This producer is not included in standard builds.
To enable it you need to trigger a custom build with native plugins enabled.
The systemd consumer allows to read from the systemd journal.


Parameters
----------

**Enable**
  Enable switches the consumer on or off.
  By default this value is set to true.

**ID**
  ID allows this producer to be found by other plugins by name.
  By default this is set to "" which does not register this producer.

**Channel**
  Channel sets the size of the channel used to communicate messages.
  By default this value is set to 8192.

**ChannelTimeoutMs**
  ChannelTimeoutMs sets a timeout in milliseconds for messages to wait if this producer's queue is full.
  A timeout of -1 or lower will drop the message without notice.
  A timeout of 0 will block until the queue is free.
  This is the default.
  A timeout of 1 or higher will wait x milliseconds for the queues to become available again.
  If this does not happen, the message will be send to the retry channel.

**ShutdownTimeoutMs**
  ShutdownTimeoutMs sets a timeout in milliseconds that will be used to detect a blocking producer during shutdown.
  By default this is set to 3 seconds.
  If processing a message takes longer to process than this duration, messages will be dropped during shutdown.

**Stream**
  Stream contains either a single string or a list of strings defining the message channels this producer will consume.
  By default this is set to "*" which means "listen to all streams but the internal".

**DropToStream**
  DropToStream defines the stream used for messages that are dropped after a timeout (see ChannelTimeoutMs).
  By default this is _DROPPED_.

**Formatter**
  Formatter sets a formatter to use.
  Each formatter has its own set of options which can be set here, too.
  By default this is set to format.Forward.
  Each producer decides if and when to use a Formatter.

**Filter**
  Filter sets a filter that is applied before formatting, i.e. before a message is send to the message queue.
  If a producer requires filtering after formatting it has to define a separate filter as the producer decides if and where to format.

**SystemdUnit**
  SystemdUnit defines what journal will be followed.
  This uses journal.add_match with _SYSTEMD_UNIT.
  By default this is set to "", which disables the filter.

**DefaultOffset**
  DefaultOffset defines where to start reading the file.
  Valid values are "oldest" and "newest".
  If OffsetFile is defined the DefaultOffset setting will be ignored unless the file does not exist.
  By default this is set to "newest".

**OffsetFile**
  OffsetFile defines the path to a file that stores the current offset.
  If the consumer is restarted that offset is used to continue reading.
  By default this is set to "" which disables the offset file.

Example
-------

.. code-block:: yaml

	- "native.Systemd":
	    Enable: true
	    ID: ""
	    Channel: 8192
	    ChannelTimeoutMs: 0
	    ShutdownTimeoutMs: 3000
	    Formatter: "format.Forward"
	    Filter: "filter.All"
	    DropToStream: "_DROPPED_"
	    Stream:
	        - "foo"
	        - "bar"
	    SystemdUnit: "sshd.service"
	    DefaultOffset: "Newest"
	    OffsetFile: ""
