<div align="center">

# WP Sec Adv

[![Test](https://github.com/typisttech/wpsecadv/actions/workflows/test.yml/badge.svg)](https://github.com/typisttech/wpsecadv/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/typisttech/wpsecadv/graph/badge.svg?token=PVY82NZYZE)](https://codecov.io/gh/typisttech/wpsecadv)
[![Go Report Card](https://goreportcard.com/badge/github.com/typisttech/wpsecadv)](https://goreportcard.com/report/github.com/typisttech/wpsecadv)
[![License](https://img.shields.io/github/license/typisttech/wpsecadv.svg)](https://github.com/typisttech/wpsecadv/blob/master/LICENSE)
[![Follow @TangRufus on X](https://img.shields.io/badge/Follow-TangRufus-15202B?logo=x&logoColor=white)](https://x.com/tangrufus)
[![Follow @TangRufus.com on Bluesky](https://img.shields.io/badge/Bluesky-TangRufus.com-blue?logo=bluesky)](https://bsky.app/profile/tangrufus.com)
[![Sponsor @TangRufus via GitHub](https://img.shields.io/badge/Sponsor-TangRufus-EA4AAA?logo=githubsponsors)](https://github.com/sponsors/tangrufus)
[![Hire Typist Tech](https://img.shields.io/badge/Hire-Typist%20Tech-778899)](https://typist.tech/contact/)

<p>
  <strong>Composer repository for WordPress security advisories.</strong>
  <br />
  <br />
  Built with ♥ by <a href="https://typist.tech/">Typist Tech</a>
</p>

</div>

---

> [!TIP]
> **Hire Tang Rufus!**
>
> I am looking for my next role, freelance or full-time.
> If you find this tool useful, I can build you more weird stuff like this.
> Let's talk if you are hiring PHP / Ruby / Go developers.
>
> Contact me at https://typist.tech/contact/

---

## Usage

```sh
composer repo --append add wpsecadv composer https://repo-wpsecadv.typist.tech
composer audit
```

You should see audit report like this:

```
Found 2 security vulnerability advisories affecting 1 package:
+-------------------+----------------------------------------------------------------------------------+
| Package           | roots/wordpress-no-content                                                       |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/112ed4f2-fe91-4d83-a3f7-eaf889870af4/wordpress                       |
| CVE               | CVE-2022-3590                                                                    |

// snip...
```

The `composer repo` subcommand is added since v2.9.0. If you are using an older Composer version, manually **append** it to your `composer.json`:

```diff
    "repositories": [
      {
        "name": "wp-composer",
        "type": "composer",
        "url": "https://repo.wp-composer.com",
        "only": [
          "wp-plugin/*",
          "wp-theme/*"
        ]
+     },
+     {
+       "name": "wpsecadv",
+       "type": "composer",
+       "url": "https://repo-wpsecadv.typist.tech"
      }
    ],
```

## Ignoring Vulnerabilities

Every WordPress core package comes with at least 2 CVE advisories.
Composer resolver blocks known vulnerabilities and fails `compsoer install|update|require`.

To ignore the advisories:

```sh
composer config audit.ignore --merge --json '["CVE-2017-14990", "CVE-2022-3590"]'
```

If you already have `audit.ignore` in object form, manually edit `composer.json`

```diff
      "config": {
          "audit": {
              "ignore": {
                  "CVE-1234": {
                      "apply": "audit",
                      "reason": "Your existing config"
+                 },
+                 "CVE-2017-14990": {
+                     "apply": "block",
+                     "reason": "XXX"
+                 },
+                 "CVE-2022-3590": {
+                     "apply": "all",
+                     "reason": "YYY"
                  }
              }
          }
      }
```

Learn more at https://getcomposer.org/doc/06-config.md#audit

## Caveats

TODO!

> [!TIP]
> **Hire Tang Rufus!**
>
> There is no need to understand any of these quirks.
> Let me handle them for you.
> I am seeking my next job, freelance or full-time.
>
> If you are hiring PHP / Ruby / Go developers,
> contact me at https://typist.tech/contact/

## Self-host

TODO!

## Credits

[`WP Sec Adv`](https://github.com/typisttech/wpsecadv) is a [Typist Tech](https://typist.tech) project and maintained by [Tang Rufus](https://x.com/TangRufus), freelance developer for [hire](https://typist.tech/contact/).

Full list of contributors can be found [here](https://github.com/typisttech/wpsecadv/graphs/contributors).

## Copyright and License

This project is a [free software](https://www.gnu.org/philosophy/free-sw.en.html) distributed under the terms of the MIT license. For the full license, see [LICENSE](./LICENSE).

## Contribute

Feedbacks / bug reports / pull requests are welcome.
