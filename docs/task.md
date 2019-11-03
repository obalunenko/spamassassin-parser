Имеется функция, которая "стягивает" имейлы с почтовых ящиков,
находит testID (string) в заголовке или теле письма.
Если емаил прошел через Spamassassin, функция вынимает репорт (см. header1, header2)
и передает данные в функцию обработки репорта. Последняя обрабатывает репорт, формирует JSON и отправляет его вместе с testID в канал.

Кто-то читает данные из канала и обновляет базу данных, находя соответствующий тест на основе testID
  

Задание:

написать модуль (модули), имеющие
1. функцию обработки репорта
2. функцию чтения из канала (просто выводить полученные данные в консоль)

На выходе должно быть : скармливаем репорт + testID, в консоли видим JSON

Пример Spamassassin репорта (2 варианта header1 и header2):



    header1 = ` * -0.0 RCVD_IN_DNSWL_NONE RBL: Sender listed at
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
`



    header2 = `Spam detection software, running on the system "server.glocksoft.com",
    has NOT identified this incoming email as spam.  The original
    message has been attached to this so you can view it or label
    similar future email.  If you have any questions, see
    root\@localhost for details.
    Content preview:  Thanks вЂ” and can I manage multiple emails from a single
      account?В  On Thu, Oct 17, 2019 at 2:55 AM, G-Lock Software < support@glocksoft.com
       > wrote: > > Hi Nick, > В  > Thank you for your reply. > > We provide the
       ability to create a custom plan for your needs. For > example, you can create
       a plan for the desired number of spam tests per > month w [...] 
    Content analysis details:   (-0.2 points, 10.0 required)
     pts rule name              description
    ---- ---------------------- --------------------------------------------------
     0.0 URIBL_BLOCKED          ADMINISTRATOR NOTICE: The query to URIBL was
                                blocked.  See
                                http://wiki.apache.org/spamassassin/DnsBlocklists#dnsbl-block
                                 for more information.
                                [URIs: glocksoft.com]
     0.0 FREEMAIL_FROM          Sender email is commonly abused enduser mail
                                provider (nickedwar[at]gmail.com)
    -0.0 SPF_PASS               SPF: sender matches SPF record
     0.0 HTML_MESSAGE           BODY: HTML included in message
    -0.1 DKIM_VALID_AU          Message has a valid DKIM or DK signature from
                                author's domain
    -0.1 DKIM_VALID             Message has at least one valid DKIM or DK signature
    -0.1 DKIM_VALID_EF          Message has a valid DKIM or DK signature from
                                envelope-from domain
     0.1 DKIM_SIGNED            Message has a DKIM or DK signature, not necessarily
                                valid
   `



//---------------------------------------------------------------------------------------------------------- 



Output:

 "spamAssassin" : {
        "score" : 1,
        "headers" : [ 
            {
                "score" : 0,
                "tag" : "SPF_PASS",
                "description" : "SPF: sender matches SPF record"
            }, 
            {
                "score" : 0,
                "tag" : "T_KAM_HTML_FONT_INVALID",
                "description" : "BODY: Test for Invalidly Named or Formatted Colors in HTML"
            }, 
            {
                "score" : 0,
                "tag" : "HTML_MESSAGE",
                "description" : "BODY: HTML included in message"
            }, 
            {
                "score" : 1.1,
                "tag" : "KAM_REALLYHUGEIMGSRC",
                "description" : "RAW: Spam with image tags with ridiculously huge http urls"
            }, 
            {
                "score" : -0.1,
                "tag" : "DKIM_VALID_AU",
                "description" : "Message has a valid DKIM or DK signature from author's domain"
            }, 
            {
                "score" : -0.1,
                "tag" : "DKIM_VALID",
                "description" : "Message has at least one valid DKIM or DK signature"
            }, 
            {
                "score" : 0.1,
                "tag" : "DKIM_SIGNED",
                "description" : "Message has a DKIM or DK signature, not necessarily valid"
            }
        ]
    }