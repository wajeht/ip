# 🌎 Ip

[![Node.js CI](https://github.com/wajeht/ip/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/wajeht/ip/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/wajeht/ip/blob/main/LICENSE) [![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/wajeht/ip)

**whatismyipaddress.com** in just a few lines of code.

If the IP address was found, the lookup method returns an object with the following structure:

```javascript
// https://ip.jaw.dev/?verbose=true
// https://ip.jaw.dev/?verbose=true&json=true
{
   ip: '127.0.0.1',               // ip address
   range: [420, 69],              // <low bound of IP block>, <high bound of IP block>
   country: 'XX',                 // 2 letter ISO-3166-1 country code
   region: 'RR',                  // Up to 3 alphanumeric variable length characters as ISO 3166-2 code
                                  // For US states this is the 2 letter state
                                  // For the United Kingdom this could be ENG as a country like “England
                                  // FIPS 10-4 subcountry code
   eu: '0',                       // 1 if the country is a member state of the European Union, 0 otherwise.
   timezone: 'Country/Zone',      // Timezone from IANA Time Zone Database
   city: "City Name",             // This is the full city name
   ll: [420, 69],                 // The latitude and longitude of the city
   metro: 420,                    // Metro code
   area: 69                       // The approximate accuracy radius (km), around the latitude and longitude
}
```

# 💻 Development

Clone the repository

```bash
$ git clone https://github.com/wajeht/ip.git
```

Copy `.env.example` to `.env`

```bash
$ cp .env.example .env
```

Install dependencies

```bash
$ npm install
```

Run development server

```bash
$ npm run dev
```

Test the application

```bash
$ npm run test
```

# 📜 License

Distributed under the MIT License © wajeht. See [LICENSE](./LICENSE) for more information.
