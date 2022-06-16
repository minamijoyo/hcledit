## master (Unreleased)

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
