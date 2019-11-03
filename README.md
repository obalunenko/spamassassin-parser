# spamassassin-parser

spamassassin-parser - a command line tool that parses spam filter reports into human readable json.

## Usage:

```bash
spamassassin-parser-cli --help
```

```bash
  -report_file string
        path to report file to process
```

##Example:
report1.txt file:

```text
 * -0.0 RCVD_IN_DNSWL_NONE RBL: Sender listed at
    *      https://www.dnswl.org/, no trust
    *      [209.85.161.101 listed in list.dnswl.org]
    * -0.0 SPF_PASS SPF: sender matches SPF record
    *  0.0 SPF_HELO_NONE SPF: HELO does not publish an SPF Record
    *  0.0 HTML_MESSAGE BODY: HTML included in message
    * -0.1 DKIM_VALID Message has at least one valid DKIM or DK signature
    *  0.1 DKIM_SIGNED Message has a DKIM or DK signature, not necessarily
    *       valid
    * -0.5 R_SB_HOSTEQIP RBL: Forward-confirmed reverse DNS (FCrDNS)
    *      succeeded
    *      [0-0=1|1=LG DACOM CORPORATION|2=6.8|3=6.8|4=2661|6=0|7=19|8=3319569|9=71889|20=mail-yw1-f101|21=google.com|22=Y|23=8.0|24=8.0|25=0|40=4.1|41=4.4|43=4.3|44=5.6|45=N|46=18|48=24|53=US|54=-97.822|55=37.751|56=1000|57=1571272183]
    *  0.0 PDS_NO_HELO_DNS High profile HELO but no A record
```

Run cli

```bash
spamassassin-parser-cli -report_file report1.txt
```

Result:

```json
{
  "spamAssassin": {
    "score": -0.5,
    "headers": [
      {
        "score": 0,
        "tag": "RCVD_IN_DNSWL_NONE",
        "description": "RBL: Sender listed at https://www.dnswl.org/, no trust [209.85.161.101 listed in list.dnswl.org]"
      },
      {
        "score": 0,
        "tag": "SPF_PASS",
        "description": "SPF: sender matches SPF record"
      },
      {
        "score": 0,
        "tag": "SPF_HELO_NONE",
        "description": "SPF: HELO does not publish an SPF Record"
      },
      {
        "score": 0,
        "tag": "HTML_MESSAGE",
        "description": "BODY: HTML included in message"
      },
      {
        "score": -0.1,
        "tag": "DKIM_VALID",
        "description": "Message has at least one valid DKIM or DK signature"
      },
      {
        "score": 0.1,
        "tag": "DKIM_SIGNED",
        "description": "Message has a DKIM or DK signature, not necessarily valid"
      },
      {
        "score": -0.5,
        "tag": "R_SB_HOSTEQIP",
        "description": "RBL: Forward-confirmed reverse DNS (FCrDNS) succeeded [0-0=1|1=LG DACOM CORPORATION|2=6.8|3=6.8|4=2661|6=0|7=19|8=3319569|9=71889|20=mail-yw1-f101|21=google.com|22=Y|23=8.0|24=8.0|25=0|40=4.1|41=4.4|43=4.3|44=5.6|45=N|46=18|48=24|53=US|54=-97.822|55=37.751|56=1000|57=1571272183]"
      },
      {
        "score": 0,
        "tag": "PDS_NO_HELO_DNS",
        "description": "High profile HELO but no A record"
      }
    ]
  }
}
```
