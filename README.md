gomaillex
=========
gomaillex is an efficient parser which allows user to parse Postfix email log files. gomaillex provides analysis in the following areas:


   * Total number of:
      * Messages received, delivered, forwarded, deferred, bounced and rejected
      * Bytes in messages received and delivered
      * Sending and Recipient Hosts/Domains
      * Senders and Recipients
      * Optional SMTPD totals for number of connections, number of hosts/domains connecting, average connect time and total connect time
   * Per-Day Traffic Summary (for multi-day logs)
   * Per-Hour Traffic (daily average for multi-day logs)
   * Optional Per-Hour and Per-Day SMTPD connection summaries
   * Sorted in descending order:
      * Recipient Hosts/Domains by message count, including:
         * Number of messages sent to recipient host/domain
         * Number of bytes in messages
         * Number of defers
         * Average delivery delay
         * Maximum delivery delay
      * Sending Hosts/Domains by message and byte count
      * Optional Hosts/Domains SMTPD connection summary
      * Senders by message count
      * Recipients by message count
      * Senders by message size
      * Recipients by message size
with an option to limit these reports to the top nn.
   * A Semi-Detailed Summary of:
      * Messages deferred
      * Messages bounced
      * Messages rejected
   * Summaries of warnings, fatal errors, and panics
   * Summary of master daemon messages
   * Optional detail of messages received, sorted by domain, then sender-in-domain, with a list of recipients-per-message.
   * Optional output of "mailq" run
