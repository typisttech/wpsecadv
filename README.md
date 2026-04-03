<div align="center">

# WP Sec Adv

[![Test](https://github.com/typisttech/wpsecadv/actions/workflows/test.yml/badge.svg)](https://github.com/typisttech/wpsecadv/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/typisttech/wpsecadv/graph/badge.svg?token=PVY82NZYZE)](https://codecov.io/gh/typisttech/wpsecadv)
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

## Quick Start

```sh
composer repo --append add wpsecadv composer https://repo-wpsecadv.typist.tech
composer audit
```

It generates audit report like this:

```
Found 2 security vulnerability advisories affecting 1 package:
+-------------------+--------------------------------------------------------------+
| Package           | roots/wordpress-no-content                                   |
| Severity          | medium                                                       |
| Advisory ID       | WPSECADV/WF/112ed4f2-fe91-4d83-a3f7-eaf889870af4/wordpress   |
| CVE               | CVE-2022-3590                                                |
// ...
```

<details>

<summary>Command "repo" is not defined.</summary>

The `composer repo` subcommand is added since Composer v2.9.0. 
If you are using an older Composer version, manually **append** it to your `composer.json`:

```diff
  "repositories": [
    {
      "name": "wp-packages",
      "type": "composer",
      "url": "https://repo.wp-packages.org"
-   }
+   },
+   {
+     "name": "wpsecadv",
+     "type": "composer",
+     "url": "https://repo-wpsecadv.typist.tech"
+   }
  ],
```

</details>

## Tutorial

First, create a fresh Bedrock project and `cd` into it:

```sh
composer create-project roots/bedrock bedrock 1.30.0
cd bedrock
```

Install some vulnerabilities:

```sh
composer require wp-theme/twentyfifteen:1.1
```

Add WP Sec Adv:

```sh
composer repo --append add wpsecadv composer https://repo-wpsecadv.typist.tech
```

Checks for security vulnerability advisories for installed packages:

```sh
composer audit
// ...
// Found 3 security vulnerability advisories affecting 2 packages
// ...
```

<details>

<summary>Full console output</summary>

```console
$ composer audit
Found 3 security vulnerability advisories affecting 2 packages:
+-------------------+----------------------------------------------------------------------------------+
| Package           | roots/wordpress-no-content                                                       |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/112ed4f2-fe91-4d83-a3f7-eaf889870af4/wordpress                       |
| CVE               | CVE-2022-3590                                                                    |
| Title             | WordPress Core - All known versions - Unauthenticated Blind Server Side Request  |
|                   | Forgery                                                                          |
|                   | ### Copyright 1999-2026 The MITRE Corporation                                    |
|                   | CVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive,        |
|                   | no-charge, royalty-free, irrevocable copyright license to reproduce, prepare     |
|                   | derivative works of, publicly display, publicly perform, sublicense, and         |
|                   | distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for    |
|                   | such purposes is authorized provided that you reproduce MITRE's copyright        |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.cve.org/Legal/TermsOfUse                                             |
|                   | ### Copyright 2012-2026 Defiant Inc.                                             |
|                   | Defiant hereby grants you a perpetual, worldwide, non-exclusive, no-charge,      |
|                   | royalty-free, irrevocable copyright license to reproduce, prepare derivative     |
|                   | works of, publicly display, publicly perform, sublicense, and distribute this    |
|                   | software vulnerability information. Any copy of the software vulnerability       |
|                   | information you make for such purposes is authorized provided that you include a |
|                   | hyperlink to this vulnerability record and reproduce Defiant's copyright         |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.wordfence.com/wordfence-intelligence-terms-and-conditions/           |
| URL               | https://www.wordfence.com/threat-intel/vulnerabilities/id/112ed4f2-fe91-4d83-a3f |
|                   | 7-eaf889870af4?source=api-prod                                                   |
| Affected versions | *                                                                                |
| Reported at       | 2022-09-06T00:00:00+00:00                                                        |
+-------------------+----------------------------------------------------------------------------------+
+-------------------+----------------------------------------------------------------------------------+
| Package           | roots/wordpress-no-content                                                       |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/9fda5e15-fdf9-4b67-93d3-2dbfa94aefe9/wordpress                       |
| CVE               | CVE-2017-14990                                                                   |
| Title             | WordPress Core - All Known Versions - Cleartext Storage of                       |
|                   | wp_signups.activation_key                                                        |
|                   | ### Copyright 1999-2026 The MITRE Corporation                                    |
|                   | CVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive,        |
|                   | no-charge, royalty-free, irrevocable copyright license to reproduce, prepare     |
|                   | derivative works of, publicly display, publicly perform, sublicense, and         |
|                   | distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for    |
|                   | such purposes is authorized provided that you reproduce MITRE's copyright        |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.cve.org/Legal/TermsOfUse                                             |
|                   | ### Copyright 2012-2026 Defiant Inc.                                             |
|                   | Defiant hereby grants you a perpetual, worldwide, non-exclusive, no-charge,      |
|                   | royalty-free, irrevocable copyright license to reproduce, prepare derivative     |
|                   | works of, publicly display, publicly perform, sublicense, and distribute this    |
|                   | software vulnerability information. Any copy of the software vulnerability       |
|                   | information you make for such purposes is authorized provided that you include a |
|                   | hyperlink to this vulnerability record and reproduce Defiant's copyright         |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.wordfence.com/wordfence-intelligence-terms-and-conditions/           |
| URL               | https://www.wordfence.com/threat-intel/vulnerabilities/id/9fda5e15-fdf9-4b67-93d |
|                   | 3-2dbfa94aefe9?source=api-prod                                                   |
| Affected versions | *                                                                                |
| Reported at       | 2017-10-10T00:00:00+00:00                                                        |
+-------------------+----------------------------------------------------------------------------------+
+-------------------+----------------------------------------------------------------------------------+
| Package           | wp-theme/twentyfifteen                                                           |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/57666105-81e4-4ef4-8889-9ce9995d2629/twentyfifteen                   |
| CVE               | CVE-2015-3429                                                                    |
| Title             | Twenty Fifteen Theme <= 1.1 & WordPress Core < 4.2.2 - Cross-Site Scripting via  |
|                   | example.html                                                                     |
|                   | ### Copyright 1999-2026 The MITRE Corporation                                    |
|                   | CVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive,        |
|                   | no-charge, royalty-free, irrevocable copyright license to reproduce, prepare     |
|                   | derivative works of, publicly display, publicly perform, sublicense, and         |
|                   | distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for    |
|                   | such purposes is authorized provided that you reproduce MITRE's copyright        |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.cve.org/Legal/TermsOfUse                                             |
|                   | ### Copyright 2012-2026 Defiant Inc.                                             |
|                   | Defiant hereby grants you a perpetual, worldwide, non-exclusive, no-charge,      |
|                   | royalty-free, irrevocable copyright license to reproduce, prepare derivative     |
|                   | works of, publicly display, publicly perform, sublicense, and distribute this    |
|                   | software vulnerability information. Any copy of the software vulnerability       |
|                   | information you make for such purposes is authorized provided that you include a |
|                   | hyperlink to this vulnerability record and reproduce Defiant's copyright         |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.wordfence.com/wordfence-intelligence-terms-and-conditions/           |
| URL               | https://www.wordfence.com/threat-intel/vulnerabilities/id/57666105-81e4-4ef4-888 |
|                   | 9-9ce9995d2629?source=api-prod                                                   |
| Affected versions | <=1.1                                                                            |
| Reported at       | 2015-04-08T00:00:00+00:00                                                        |
+-------------------+----------------------------------------------------------------------------------+
```

</details>


The best course of action is to update packages to patched versions.

Update the Twenty Fifteen theme:

```sh
composer require wp-theme/twentyfifteen
// ...
// Found 2 security vulnerability advisories affecting 1 package
// ...
```

<details>

<summary>Full console output</summary>

```console
$ composer require wp-theme/twentyfifteen
./composer.json has been updated
Running composer update wp-theme/twentyfifteen
Loading composer repositories with package information
Updating dependencies
Lock file operations: 0 installs, 1 update, 0 removals
  - Upgrading wp-theme/twentyfifteen (1.1 => 4.1)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 0 installs, 1 update, 0 removals
  - Upgrading wp-theme/twentyfifteen (1.1 => 4.1): Extracting archive
Generating optimized autoload files
Found 2 security vulnerability advisories affecting 1 package.
Run "composer audit" for a full list of advisories.
Using version ^4.1 for wp-theme/twentyfifteen
```

</details>

However, there may not be a patch yet or never will be (as the two WordPress core CVEs).

> [!WARNING]
> Blindly ignoring packages from secutiy blockings is **dangerous**.
>
> You should do so only in exceptional cases.

Ignore `roots/wordpress-no-content` from auditing, edit `composer.json`:

```json
{
  "config": {
    "audit": {
      "ignore": ["roots/wordpress-no-content"]
    }
  }
}
```

When installing packages with known vulnerabilities, Composer resolver blocks them and fails `composer update|require`.

Install a vulnerable WooCommerce version:

```sh
composer require wp-plugin/woocommerce:10.5.0
// ...
// Your requirements could not be resolved to an installable set of packages.
//
//  Problem 1
//    - Root composer.json requires wp-plugin/woocommerce 10.5.0 (exact version match: 10.5.0 or 10.5.0.0), found wp-plugin/woocommerce[10.5.0] but these were not loaded, because they are affected by security advisories ("WPSECADV/WF/df7eca9b-e353-49e7-8706-89c1787637e9/woocommerce").
// ...
// Installation failed, reverting ./composer.json and ./composer.lock to their original content.
```

<details>

<summary>Full console output</summary>

```console
$ composer require wp-plugin/woocommerce:10.5.0
./composer.json has been updated
Running composer update wp-plugin/woocommerce
Loading composer repositories with package information
Updating dependencies
Your requirements could not be resolved to an installable set of packages.

  Problem 1
    - Root composer.json requires wp-plugin/woocommerce 10.5.0 (exact version match: 10.5.0 or 10.5.0.0), found wp-plugin/woocommerce[10.5.0] but these were not loaded, because they are affected by security advisories ("WPSECADV/WF/df7eca9b-e353-49e7-8706-89c1787637e9/woocommerce"). Go to https://packagist.org/security-advisories/ to find advisory details. To ignore the advisories, add them to the audit "ignore" config. To turn the feature off entirely, you can set "block-insecure" to false in your "audit" config.


Installation failed, reverting ./composer.json and ./composer.lock to their original content.
```

</details>

Unfortunately, a WooCommerce add-on compatibility issue forces us to stay with WooCommerce v10.5.0.

To disable security blocking during install:

```sh
composer require wp-plugin/woocommerce:10.5.0 --no-security-blocking
// ...
// Found 2 ignored security vulnerability advisories affecting 1 package.
// Found 1 security vulnerability advisory affecting 1 package.
// ...
```

<details>

<summary>Full console output</summary>

```console
$ composer require wp-plugin/woocommerce:10.5.0 --no-security-blocking
./composer.json has been updated
Running composer update wp-plugin/woocommerce
Loading composer repositories with package information
Updating dependencies
Lock file operations: 1 install, 0 updates, 0 removals
  - Locking wp-plugin/woocommerce (10.5.0)
Writing lock file
Installing dependencies from lock file (including require-dev)
Package operations: 1 install, 0 updates, 0 removals
  - Installing wp-plugin/woocommerce (10.5.0): Extracting archive
Generating optimized autoload files
Found 2 ignored security vulnerability advisories affecting 1 package.
Found 1 security vulnerability advisory affecting 1 package.
Run "composer audit" for a full list of advisories.
```

</details>

The `--no-security-blocking` flag allows installing packages with security advisories but it is one-off.
Future `composer update|require` will be blocked.

Once you have it installed, get the CVE IDs via:

```sh
composer audit
// ...
// | Package   | wp-plugin/woocommerce   |
// | CVE       | CVE-2026-3589           |
// ...
```

<details>

<summary>Full console output</summary>

```console
$ composer audit
Found 2 ignored security vulnerability advisories affecting 1 package:
+-------------------+----------------------------------------------------------------------------------+
| Package           | roots/wordpress-no-content                                                       |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/112ed4f2-fe91-4d83-a3f7-eaf889870af4/wordpress                       |
| CVE               | CVE-2022-3590                                                                    |
| Title             | WordPress Core - All known versions - Unauthenticated Blind Server Side Request  |
|                   | Forgery                                                                          |
|                   | ### Copyright 1999-2026 The MITRE Corporation                                    |
|                   | CVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive,        |
|                   | no-charge, royalty-free, irrevocable copyright license to reproduce, prepare     |
|                   | derivative works of, publicly display, publicly perform, sublicense, and         |
|                   | distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for    |
|                   | such purposes is authorized provided that you reproduce MITRE's copyright        |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.cve.org/Legal/TermsOfUse                                             |
|                   | ### Copyright 2012-2026 Defiant Inc.                                             |
|                   | Defiant hereby grants you a perpetual, worldwide, non-exclusive, no-charge,      |
|                   | royalty-free, irrevocable copyright license to reproduce, prepare derivative     |
|                   | works of, publicly display, publicly perform, sublicense, and distribute this    |
|                   | software vulnerability information. Any copy of the software vulnerability       |
|                   | information you make for such purposes is authorized provided that you include a |
|                   | hyperlink to this vulnerability record and reproduce Defiant's copyright         |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.wordfence.com/wordfence-intelligence-terms-and-conditions/           |
| URL               | https://www.wordfence.com/threat-intel/vulnerabilities/id/112ed4f2-fe91-4d83-a3f |
|                   | 7-eaf889870af4?source=api-prod                                                   |
| Affected versions | *                                                                                |
| Reported at       | 2022-09-06T00:00:00+00:00                                                        |
| Ignore reason     | None specified                                                                   |
+-------------------+----------------------------------------------------------------------------------+
+-------------------+----------------------------------------------------------------------------------+
| Package           | roots/wordpress-no-content                                                       |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/9fda5e15-fdf9-4b67-93d3-2dbfa94aefe9/wordpress                       |
| CVE               | CVE-2017-14990                                                                   |
| Title             | WordPress Core - All Known Versions - Cleartext Storage of                       |
|                   | wp_signups.activation_key                                                        |
|                   | ### Copyright 1999-2026 The MITRE Corporation                                    |
|                   | CVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive,        |
|                   | no-charge, royalty-free, irrevocable copyright license to reproduce, prepare     |
|                   | derivative works of, publicly display, publicly perform, sublicense, and         |
|                   | distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for    |
|                   | such purposes is authorized provided that you reproduce MITRE's copyright        |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.cve.org/Legal/TermsOfUse                                             |
|                   | ### Copyright 2012-2026 Defiant Inc.                                             |
|                   | Defiant hereby grants you a perpetual, worldwide, non-exclusive, no-charge,      |
|                   | royalty-free, irrevocable copyright license to reproduce, prepare derivative     |
|                   | works of, publicly display, publicly perform, sublicense, and distribute this    |
|                   | software vulnerability information. Any copy of the software vulnerability       |
|                   | information you make for such purposes is authorized provided that you include a |
|                   | hyperlink to this vulnerability record and reproduce Defiant's copyright         |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.wordfence.com/wordfence-intelligence-terms-and-conditions/           |
| URL               | https://www.wordfence.com/threat-intel/vulnerabilities/id/9fda5e15-fdf9-4b67-93d |
|                   | 3-2dbfa94aefe9?source=api-prod                                                   |
| Affected versions | *                                                                                |
| Reported at       | 2017-10-10T00:00:00+00:00                                                        |
| Ignore reason     | None specified                                                                   |
+-------------------+----------------------------------------------------------------------------------+
Found 1 security vulnerability advisory affecting 1 package:
+-------------------+----------------------------------------------------------------------------------+
| Package           | wp-plugin/woocommerce                                                            |
| Severity          | medium                                                                           |
| Advisory ID       | WPSECADV/WF/df7eca9b-e353-49e7-8706-89c1787637e9/woocommerce                     |
| CVE               | CVE-2026-3589                                                                    |
| Title             | WooCommerce < 10.5.3 - Cross-Site Request Forgery                                |
|                   | ### Copyright 1999-2026 The MITRE Corporation                                    |
|                   | CVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive,        |
|                   | no-charge, royalty-free, irrevocable copyright license to reproduce, prepare     |
|                   | derivative works of, publicly display, publicly perform, sublicense, and         |
|                   | distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for    |
|                   | such purposes is authorized provided that you reproduce MITRE's copyright        |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.cve.org/Legal/TermsOfUse                                             |
|                   | ### Copyright 2012-2026 Defiant Inc.                                             |
|                   | Defiant hereby grants you a perpetual, worldwide, non-exclusive, no-charge,      |
|                   | royalty-free, irrevocable copyright license to reproduce, prepare derivative     |
|                   | works of, publicly display, publicly perform, sublicense, and distribute this    |
|                   | software vulnerability information. Any copy of the software vulnerability       |
|                   | information you make for such purposes is authorized provided that you include a |
|                   | hyperlink to this vulnerability record and reproduce Defiant's copyright         |
|                   | designation and this license in any such copy.                                   |
|                   | https://www.wordfence.com/wordfence-intelligence-terms-and-conditions/           |
| URL               | https://www.wordfence.com/threat-intel/vulnerabilities/id/df7eca9b-e353-49e7-870 |
|                   | 6-89c1787637e9?source=api-prod                                                   |
| Affected versions | <10.5.3                                                                          |
| Reported at       | 2026-03-10T00:00:00+00:00                                                        |
+-------------------+----------------------------------------------------------------------------------+
```

</details>

Allow specific advisories to be installed, edit `composer.json`:

```json
{
  "config": {
    "audit": {
      "ignore": {
        "roots/wordpress-no-content": {
          "apply": "all",
          "reason": "We live dangerously and don't care about this one"
        },
        "CVE-2026-3589": {
          "apply": "block",
          "reason": "Waiting for FooBar add-on v1.2.3 to be released. Allow during updates but still report in audits"
        }
      }
    }
  }
}
```

All of the above are Composer features. WP Sec Adv merely makes Wordfence vulnerability data feed available in Composer consumable format.

Learn more at:

- https://getcomposer.org/doc/06-config.md#audit
- https://getcomposer.org/doc/03-cli.md#audit
- https://blog.packagist.com/discover-security-advisories-with-composers-audit-command/
- https://www.wordfence.com/help/wordfence-intelligence/v3-accessing-and-consuming-the-vulnerability-data-feed/

> [!TIP]
> **Hire Tang Rufus!**
>
> There is no need to understand any of these quirks.
> Let me handle them for you.
> I am seeking my next job, freelance or full-time.
>
> If you are hiring PHP / Ruby / Go developers,
> contact me at https://typist.tech/contact/

## Disable Security Blocking

Besides the one-off `--no-security-blocking` flag, you can persistently disable security blocking by:

```sh
composer config audit.block-insecure false
```

Or, manually edit `composer.json`:

```json
{
  "config": {
    "audit": {
      "block-insecure": false
    }
  }
}
```

## Package Resolving

Composer package names consist of `vendor` and `project`, e.g: `my-vendor/my-project` whereas WordPress themes and plugins are identified by `slug` only.

WP Sec Adv matches Composer packages with WordPress themes & plugins by `project` and `slug`. For example:

| Composer                            | WordPress          |
| ----------------------------------- | ------------------ |
| `wp-plugin/woocommerce`             | `woocommerce`      |
| `wpackagist-plugin/woocommerce`     | `woocommerce`      |
| `my-mirror/woocommerce`             | `woocommerce`      |
| `gravity/gravityforms`              | `gravityforms`     |
| `my-mirror/gravityforms`            | `gravityforms`     |
| `wp-theme/twentytwentyfive`         | `twentytwentyfive` |
| `wpackagist-theme/twentytwentyfive` | `twentytwentyfive` |
| `my-mirror/twentytwentyfive`        | `twentytwentyfive` |

### `exclude`

In case of naming collision, add `exclude` to the repository config.

For example, this setup prevents mismatching `spatie/ignition` as the [Ignition theme](https://www.wordfence.com/threat-intel/vulnerabilities/wordpress-themes/ignition#:~:text=Ignition):

```diff
  "repositories": [
    {
      "name": "wp-packages",
      "type": "composer",
      "url": "https://repo.wp-packages.org"
    },
    {
      "name": "wpsecadv",
      "type": "composer",
-       "url": "https://repo-wpsecadv.typist.tech"
+       "url": "https://repo-wpsecadv.typist.tech",
+       "exclude": [
+         "spatie/ignition"
+       ]
    }
  ],
```

### `only`

To avoid mismatches and speed up Composer operations, add `only` to the repository config:

```diff
  "repositories": [
    {
      "name": "wp-packages",
      "type": "composer",
      "url": "https://repo.wp-packages.org"
    },
    {
      "name": "wpsecadv",
      "type": "composer",
-     "url": "https://repo-wpsecadv.typist.tech"
+     "url": "https://repo-wpsecadv.typist.tech",
+     "only": [
+       "wp-plugin/*",
+       "wp-theme/*",
+       "wp-core/*",
+       "wpackagist-plugin/*",
+       "wpackagist-theme/*",
+       "roots/wordpress-no-content",
+       "roots/wordpress-full",
+       "johnpbloch/wordpress-core",
+       "deliciousbrains-plugin/*",
+       "gravity/*",
+       "yoast/*",
+       "my-mirror/*"
+     ]
    }
  ],
```

Adjust the `only` array to suit your situation.

## Continuous Monitoring

> [!IMPORTANT]
> Vulnerabilities get discovered every day. Audit your dependencies **automatically**.

### GitHub Actions

```yml
name: Audit Dependencies

on:
  workflow_dispatch:
  schedule:
    - cron: '0 9 * * *' # Once a day
  pull_request:
  push:

permissions:
  contents: read

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout composer.json & composer.lock
        uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2
        with:
          persist-credentials: false
          sparse-checkout: |
            composer.json
            composer.lock

      - name: Setup PHP
        uses: shivammathur/setup-php@accd6127cb78bee3e8082180cb391013d204ef9f # v2.37.0
        with:
          php-version: '8.5'

      - name: Checks for security vulnerability advisories
        run: composer audit --locked
```

## Best Practices

- Prefer the detailed `config.audit.ignore` object with [`apply` and `reason`](https://getcomposer.org/doc/06-config.md#detailed-format-with-apply-scope-) so you can review the decisions in the future
- Unless you have [continuous monitoring](#continuous-monitoring) set up, use [`config.audit.block-insecure`](https://getcomposer.org/doc/06-config.md#block-insecure) only as an emergency or short-term measure
- Narrow [`only`](#only) just enough to cover your WordPress core, plugins and themes
- Check the vulnerability advisory details. Even though it gets patched, the damage might already be done

## Self-host

### Fly.io

First, install `fly` (aka `flyctl`) if you haven't, learn more at https://fly.io/docs/flyctl/install/

Then, launch a new app via [fly.toml](fly.toml):

```sh
# Download fly.toml from GitHub
curl -O https://raw.githubusercontent.com/typisttech/wpsecadv/refs/heads/main/fly.toml

# Launch the App
fly launch --copy-config

# Verify
fly apps open /packages.json
```

To update advisory data, re-deploy the [`latest`](https://github.com/typisttech/wpsecadv/pkgs/container/wpsecadv/latest) container image:

```sh
fly deploy
```

To auto-deploy via GitHub Actions, see [`deploy.yml`](.github/workflows/deploy.yml).

## Wordfence

WP Sec Adv sources the advisory data from [Wordfence vulnerability data feed](https://www.wordfence.com/help/wordfence-intelligence/v3-accessing-and-consuming-the-vulnerability-data-feed/). Kudos to the Wordfence team for opening the data feed freely to all.

The data feed comes with [attribution requirement](https://www.wordfence.com/help/wordfence-intelligence/v3-accessing-and-consuming-the-vulnerability-data-feed/#mitre_attribution_requirement). However, Composer has no mechanism to display the copyrights. Thus, WP Sec Adv appends copyright details to advisory titles.

## Credits

[`WP Sec Adv`](https://github.com/typisttech/wpsecadv) is a [Typist Tech](https://typist.tech) project and maintained by [Tang Rufus](https://x.com/TangRufus), freelance developer for [hire](https://typist.tech/contact/).

Full list of contributors can be found [here](https://github.com/typisttech/wpsecadv/graphs/contributors).

## Copyright and License

This project is a [free software](https://www.gnu.org/philosophy/free-sw.en.html) distributed under the terms of the MIT license. For the full license, see [LICENSE](./LICENSE).

## Contribute

Feedbacks / bug reports / pull requests are welcome.
