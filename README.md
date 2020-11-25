# spamassassin-parser

[![GO](https://img.shields.io/github/go-mod/go-version/oleg-balunenko/spamassassin-parser)](https://golang.org/doc/devel/release.html)
[![Build Status](https://travis-ci.com/obalunenko/spamassassin-parser.svg?branch=master)](https://travis-ci.com/obalunenko/spamassassin-parser)
[![Coverage Status](https://coveralls.io/repos/github/obalunenko/spamassassin-parser/badge.svg?branch=master)](https://coveralls.io/github/obalunenko/spamassassin-parser?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/obalunenko/spamassassin-parser)](https://goreportcard.com/report/github.com/obalunenko/spamassassin-parser)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=oleg-balunenko_spamassassin-parser&metric=alert_status)](https://sonarcloud.io/dashboard?id=oleg-balunenko_spamassassin-parser)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/8847ad100b3f415fa419430a58de1a2d)](https://www.codacy.com/manual/oleg.balunenko/spamassassin-parser?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=oleg-balunenko/spamassassin-parser&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/obalunenko/spamassassin-parser?status.svg)](https://godoc.org/github.com/obalunenko/spamassassin-parser)
[![Latest release artifacts](https://img.shields.io/github/v/release/obalunenko/spamassassin-parser)](https://github.com/obalunenko/spamassassin-parser/releases/latest)
[![Docker pulls](https://img.shields.io/docker/pulls/olegbalunenko/spamassassin-parser)](https://hub.docker.com/r/olegbalunenko/spamassassin-parser)
[![License](https://img.shields.io/github/license/obalunenko/spamassassin-parser)](/LICENSE)

<p align="center">
  <img src="https://github.com/obalunenko/spamassassin-parser/blob/master/.assets/assassingopher.png" alt="" width="300">
  <br>
</p>

spamassassin-parser - a command line tool that parses spam filter reports into human readable json.

## Usage

1. Download executable file: [![Latest release artifacts](https://img.shields.io/github/v/release/obalunenko/spamassassin-parser)](https://github.com/oleg-balunenko/spamassassin-parser/releases/latest)
2. Unrar archive.
3. a. Run executable `spamassassin-parser`
   b. Run docker-compose `docker-compose -f ./docker-compose.yml up --build -d`

Environment variables used:

```env
  SPAMASSASSIN_INPUT: Path to directory where files for proccession are located (default "input")
  SPAMASSASSIN_RESULT: Path to directory where parserd results will be stored (default "result")
  SPAMASSASSIN_ARCHIVE: Path to dir where processed files will be moved for history (default "archive")
  SPAMASSASSIN_RECEIVE_ERRORS: Boolean value to enable receive errors from processor, if false - will be   just logged (default: "true")
```

## Example

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

Run service

```bash
spamassassin-parser
```

- Now application will poll the directory input for new files with txt extension.
- Put a new file for procession.
- After a file processed parsed result will be stored in the file at output directory and original file will be moved to archive.

Result example:

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
