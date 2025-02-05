## master (Unreleased)

## 0.2.16 (2025/02/05)

NEW FEATURES:

* Add block new cmd ([#106](https://github.com/minamijoyo/hcledit/pull/106)) ([#109](https://github.com/minamijoyo/hcledit/pull/109))  ([#110](https://github.com/minamijoyo/hcledit/pull/110))

ENHANCEMENTS:

* Update Go to v1.23.0 ([#107](https://github.com/minamijoyo/hcledit/pull/107))
* Update hcl to v2.23.0 ([#108](https://github.com/minamijoyo/hcledit/pull/108))

## 0.2.15 (2024/08/30)

BUG FIXES:

* Fix attribute get --with-comments for inline comments ([#104](https://github.com/minamijoyo/hcledit/pull/104))

## 0.2.14 (2024/08/29)

NEW FEATURES:

* add --with-comments to preserve comments when returning an attribute ([#103](https://github.com/minamijoyo/hcledit/pull/103))

## 0.2.13 (2024/08/02)

BUG FIXES:

* Fix syntax error in release workflow ([#101](https://github.com/minamijoyo/hcledit/pull/101))

## 0.2.12 (2024/08/02)

ENHANCEMENTS:

* Update hcl to v2.21.0 ([#95](https://github.com/minamijoyo/hcledit/pull/95))
* Update alpine to v3.20 ([#96](https://github.com/minamijoyo/hcledit/pull/96))
* Update golangci lint to v1.59.1 ([#97](https://github.com/minamijoyo/hcledit/pull/97))
* Update setup-go to v5 ([#98](https://github.com/minamijoyo/hcledit/pull/98))
* Update goreleaser to v2 ([#99](https://github.com/minamijoyo/hcledit/pull/99))
* Switch to the official action for creating GitHub App token ([#100](https://github.com/minamijoyo/hcledit/pull/100))

## 0.2.11 (2024/04/15)

ENHANCEMENTS:

* feat: update to use go 1.22 ([#91](https://github.com/minamijoyo/hcledit/pull/91))
* Add support for namespaced function ([#93](https://github.com/minamijoyo/hcledit/pull/93))

## 0.2.10 (2023/09/20)

NEW FEATURES:

* feat: add support for escaping . in address ([#83](https://github.com/minamijoyo/hcledit/pull/83))

ENHANCEMENTS:

* Update Go to v1.21 ([#86](https://github.com/minamijoyo/hcledit/pull/86))
* Update hcl to v2.18.0 ([#87](https://github.com/minamijoyo/hcledit/pull/87))

## 0.2.9 (2023/06/12)

ENHANCEMENTS:

* Update hcl to v2.17.0 ([#81](https://github.com/minamijoyo/hcledit/pull/81))

BUG FIXES:

* Fix unexpected format when files do not end with newline ([#79](https://github.com/minamijoyo/hcledit/pull/79))

## 0.2.8 (2023/05/11)

BUG FIXES:

* Fix multiline comment parsing ([#78](https://github.com/minamijoyo/hcledit/pull/78))

## 0.2.7 (2023/04/19)

ENHANCEMENTS:

* Update Go to v1.20 ([#73](https://github.com/minamijoyo/hcledit/pull/73))
* Update hcl to v2.16.2 ([#74](https://github.com/minamijoyo/hcledit/pull/74))
* Use a native cache feature in actions/setup-go ([#75](https://github.com/minamijoyo/hcledit/pull/75))
* Update actions/setup-go to v4 ([#76](https://github.com/minamijoyo/hcledit/pull/76))
* Add windows build ([#77](https://github.com/minamijoyo/hcledit/pull/77))

## 0.2.6 (2022/08/12)

ENHANCEMENTS:

* Use GitHub App token for updating brew formula on release ([#59](https://github.com/minamijoyo/hcledit/pull/59))

## 0.2.5 (2022/06/16)

ENHANCEMENTS:

* Update Go to v1.18.3 ([#57](https://github.com/minamijoyo/hcledit/pull/57))

## 0.2.4 (2022/06/13)

ENHANCEMENTS:

* Expose VerticalFormat ([#43](https://github.com/minamijoyo/hcledit/pull/43))
* Expose GetAttributeValueAsString ([#47](https://github.com/minamijoyo/hcledit/pull/47))
* Update golangci-lint to v1.45.2 and actions to latest ([#49](https://github.com/minamijoyo/hcledit/pull/49))
* Read Go version from .go-version on GitHub Actions ([#53](https://github.com/minamijoyo/hcledit/pull/53))
* Update Go to v1.17.10 and Alpine to v3.16 ([#54](https://github.com/minamijoyo/hcledit/pull/54))
* Update hcl to v2.12.0 ([#55](https://github.com/minamijoyo/hcledit/pull/55))

BUG FIXES:

* Trim trailing duplicated TokenNewline in VerticalFormat ([#48](https://github.com/minamijoyo/hcledit/pull/48))

## 0.2.3 (2022/02/12)

ENHANCEMENTS:

* Use golangci-lint instead of golint ([#40](https://github.com/minamijoyo/hcledit/pull/40))
* Fix lint errors ([#41](https://github.com/minamijoyo/hcledit/pull/41))
* Update hcl to v2.11.1 ([#42](https://github.com/minamijoyo/hcledit/pull/42))

## 0.2.2 (2021/11/28)

ENHANCEMENTS:

* Update Go to v1.17.3 and Alpine to 3.14 ([#38](https://github.com/minamijoyo/hcledit/pull/38))
* Update hcl to v2.10.1 ([#39](https://github.com/minamijoyo/hcledit/pull/39))
* Add Apple Silicon (ARM 64) build ([#36](https://github.com/minamijoyo/hcledit/pull/36))

## 0.2.1 (2021/10/28)

ENHANCEMENTS:

* Restrict permissions for GitHub Actions ([#34](https://github.com/minamijoyo/hcledit/pull/34))
* Set timeout for GitHub Actions ([#35](https://github.com/minamijoyo/hcledit/pull/35))

## 0.2.0 (2021/04/06)

BREAKING CHANGES:

* Skip formatter if filter didn't change contents ([#24](https://github.com/minamijoyo/hcledit/pull/24))

Previously outputs are always formatted, but the outputs are no longer formatted if a given address doesn't match to suppress meaningless diff.

NEW FEATURES:

* Add support for getting nested block ([#22](https://github.com/minamijoyo/hcledit/pull/22))
* Add body get command ([#23](https://github.com/minamijoyo/hcledit/pull/23))
* Add support for in-place update ([#25](https://github.com/minamijoyo/hcledit/pull/25))

ENHANCEMENTS:

* Redesign interfaces in editor package ([#18](https://github.com/minamijoyo/hcledit/pull/18))
* Update Go to v1.16.0 ([#19](https://github.com/minamijoyo/hcledit/pull/19))
* Update hcl to v2.9.0 ([#20](https://github.com/minamijoyo/hcledit/pull/20))

## 0.1.3 (2021/01/30)

ENHANCEMENTS:

* Update hcl to v2.8.2 ([#16](https://github.com/minamijoyo/hcledit/pull/16))
* Fix broken GitHub Actions ([#17](https://github.com/minamijoyo/hcledit/pull/17))

## 0.1.2 (2020/10/28)

NEW FEATURES:

* Add attribute append command ([#14](https://github.com/minamijoyo/hcledit/pull/14))
* Add fmt command ([#15](https://github.com/minamijoyo/hcledit/pull/15))

## 0.1.1 (2020/10/25)

NEW FEATURES:

* Add block append command ([#8](https://github.com/minamijoyo/hcledit/pull/8))

ENHANCEMENTS:

* Add integration test ([#5](https://github.com/minamijoyo/hcledit/pull/5))
* Update hcl to v2.7.0 ([#6](https://github.com/minamijoyo/hcledit/pull/6))
* Update Go to v1.15.2 ([#7](https://github.com/minamijoyo/hcledit/pull/7))
* Refactor to test argument flags ([#9](https://github.com/minamijoyo/hcledit/pull/9))
* Prevent uploading pre-release to Homebrew ([#12](https://github.com/minamijoyo/hcledit/pull/12))

BUG FIXES:

* Fix binary compatibility issue for alpine ([#11](https://github.com/minamijoyo/hcledit/pull/11))

## 0.1.0 (2020/08/22)

Initial release
