# Dependencies Version Tracking

This document tracks the version dependencies of ChronoGo's external libraries for maintenance and update tracking purposes.

## Holiday Data Dependencies

### GoHoliday Library
- **Current Version**: v0.5.3
- **Repository**: https://github.com/coredds/GoHoliday
- **Upstream Source**: Based on [Vacanza holidays](https://github.com/vacanza/holidays) framework
- **Vacanza Version Tracked**: v0.80+ (September 2025)
- **Last Updated**: August 2025
- **Supported Countries**: 33 countries with 500+ regional subdivisions

### Supported Countries (33)
US, GB, CA, AU, NZ, DE, FR, JP, IN, BR, MX, IT, ES, NL, KR, PT, PL, RU, CN, TH, SG, MY, ID, PH, VN, TW, HK, ZA, EG, NG, KE, GH, MA, TN

### Update Procedure
1. Check Vacanza holidays repository for latest releases: https://github.com/vacanza/holidays/releases
2. Review CHANGES.md for updates to supported countries: https://github.com/vacanza/holidays/blob/dev/CHANGES.md
3. Update GoHoliday dependency if new Vacanza version includes relevant fixes
4. Run comprehensive tests: `go test -v ./... -run TestAllSupportedCountries`
5. Update this tracking document with new version information

### Notes
- GoHoliday provides Go bindings for Python's Vacanza holidays framework
- Vacanza is the modern, actively maintained fork of the original python-holidays library
- Updates should focus on the 33 countries currently supported by GoHoliday
- Test all countries after updates to ensure compatibility

## Go Dependencies
- **Go Version**: 1.23+
- **Main Dependency**: github.com/coredds/GoHoliday v0.5.3

## Expansion Opportunities
- **Available for Addition**: 24+ additional countries verified as available in GoHoliday v0.5.3+
- **Priority Targets**: Switzerland (CH), Norway (NO), Sweden (SE), Israel (IL), Turkey (TR), UAE (AE), Saudi Arabia (SA)
- **Strategic Value**: Financial centers, major economies, regional completion opportunities
- **See**: EXPANSION_ROADMAP.md for detailed implementation strategy

## Last Update Check
- **Date**: January 2025
- **Vacanza Version Checked**: v0.80 (latest as of September 2025)
- **Status**: Up to date - no critical updates needed for supported countries
- **Expansion Status**: 24+ countries verified as ready for implementation
- **Next Check Recommended**: March 2025
