import * as galleta from "../dist/index";


const sampleCookie: string = "eyJTdGF0ZVN0cmluZyI6IlRlc3RWYWx1ZSIsIlNpZyI6eyJGb3JtYXQiOiJzc2gtcnNhIiwiQmxvYiI6IkNCUHdES015RThUa3hRdm90d2xCc0ZxUUx3SG53cU1IMnNCU29jTXBzRVBUbFJTMXpOZ3hvNlloTjBTQVgxL1FaRGxoblFqbDJTOXluc21MRWJFa0NHTkJpM2lMN1FKK0F4NFZtSzZuTnBXcU1wSTd6cGFIZWVMY05ZcCt5RGFRd1dFT2pjbWk1d0ExdURvN3lLclZoWlFQamVTRHVBK1ZHM0JaSjZ3OWIwb0M0di9sVDFpMDgzSmZHUWpCWW5pS1lMaGZDeC9zQko4T2xuQWgxSDZmdi9nZU1MbWhuWVpJcktDdC84K01BSm5YTWM2UEZ3TlFqVHhWOVoySUZBcDc5S05VYzl0YmtxcjRwRzFTaWszbW9rczJES1U5SjdnYVlITGxsRFB0bloxaWJKeDdUY290S2FJQmRrTHBnTkJMSkNJRjZpWnc1TGNoTlAvS1Z2MktQQT09IiwiUmVzdCI6bnVsbH19"

describe("GetCookie Tests", () => {
  it("should return a json object", () => {
    console.log(galleta.decodeCookie(sampleCookie));
    true;
  });
});