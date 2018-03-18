package main

const apiNoTags = `
[

]
`

const apiNoReleases = `
[

]
`

const apiTagsPage1 = `
[
  {
    "name": "tag-jin-tries-arbitrary-tag-01",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/tag-jin-tries-arbitrary-tag-01",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/tag-jin-tries-arbitrary-tag-01",
    "commit": {
      "sha": "fc955c9b287ebfbce5824e9b7712ab11b9d0785c",
      "url": "https://api.github.com/repos/foo/bar/commits/fc955c9b287ebfbce5824e9b7712ab11b9d0785c"
    }
  },
  {
    "name": "tag-non-semantic-version-1.0.05.6",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/tag-jin-new-docker-jenkins-test",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/tag-jin-new-docker-jenkins-test",
    "commit": {
      "sha": "389e5f670c18f093bc0d016239507401053e3962",
      "url": "https://api.github.com/repos/foo/bar/commits/389e5f670c18f093bc0d016239507401053e3962"
    }
  },
  {
    "name": "tag-jin-prefix-0.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/tag-jin-dodgy-jenkins-trigger-test",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/tag-jin-dodgy-jenkins-trigger-test",
    "commit": {
      "sha": "389e5f670c18f093bc0d016239507401053e3962",
      "url": "https://api.github.com/repos/foo/bar/commits/389e5f670c18f093bc0d016239507401053e3962"
    }
  },
  {
    "name": "tag-fix-hcc-iframe-for-ssl",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/tag-fix-hcc-iframe-for-ssl",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/tag-fix-hcc-iframe-for-ssl",
    "commit": {
      "sha": "e4aac3ef624fdfe4c169b8eeaef197b68b11366e",
      "url": "https://api.github.com/repos/foo/bar/commits/e4aac3ef624fdfe4c169b8eeaef197b68b11366e"
    }
  },
  {
    "name": "tag-fix-cid-generation",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/tag-fix-cid-generation",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/tag-fix-cid-generation",
    "commit": {
      "sha": "35c7d6bab9bddeaef07efe187ffefe10b2b8dd76",
      "url": "https://api.github.com/repos/foo/bar/commits/35c7d6bab9bddeaef07efe187ffefe10b2b8dd76"
    }
  },
  {
    "name": "foo-bar-Symbionts",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts",
    "commit": {
      "sha": "e692208467db213139b24734715111119269b92b",
      "url": "https://api.github.com/repos/foo/bar/commits/e692208467db213139b24734715111119269b92b"
    }
  },
  {
    "name": "foo-bar-Symbionts-12",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-12",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-12",
    "commit": {
      "sha": "a97dec23545aaa8528ac8ae89c130252125df5d5",
      "url": "https://api.github.com/repos/foo/bar/commits/a97dec23545aaa8528ac8ae89c130252125df5d5"
    }
  },
  {
    "name": "foo-bar-Symbionts-11",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-11",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-11",
    "commit": {
      "sha": "69e10c0603b52adf8e32eda4dd7088cb58668cea",
      "url": "https://api.github.com/repos/foo/bar/commits/69e10c0603b52adf8e32eda4dd7088cb58668cea"
    }
  },
  {
    "name": "foo-bar-Symbionts-10",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-10",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-10",
    "commit": {
      "sha": "ee1677cb631298595dd405bca5a7e70f577a6fb4",
      "url": "https://api.github.com/repos/foo/bar/commits/ee1677cb631298595dd405bca5a7e70f577a6fb4"
    }
  },
  {
    "name": "foo-bar-Symbionts-9",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-9",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-9",
    "commit": {
      "sha": "e5354477371d53af35f6797b7d0e7645c421510e",
      "url": "https://api.github.com/repos/foo/bar/commits/e5354477371d53af35f6797b7d0e7645c421510e"
    }
  },
  {
    "name": "foo-bar-Symbionts-8",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-8",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-8",
    "commit": {
      "sha": "114dac458e19f9eeb83748b399e2a6c93893c639",
      "url": "https://api.github.com/repos/foo/bar/commits/114dac458e19f9eeb83748b399e2a6c93893c639"
    }
  },
  {
    "name": "foo-bar-Symbionts-7",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-7",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-7",
    "commit": {
      "sha": "f0a58665f8964d138ff11614eeedb3f1fa99127f",
      "url": "https://api.github.com/repos/foo/bar/commits/f0a58665f8964d138ff11614eeedb3f1fa99127f"
    }
  },
  {
    "name": "foo-bar-Symbionts-6",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-6",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-6",
    "commit": {
      "sha": "28e02ce2f542e6f3170109e2e0459a1600d999b6",
      "url": "https://api.github.com/repos/foo/bar/commits/28e02ce2f542e6f3170109e2e0459a1600d999b6"
    }
  },
  {
    "name": "foo-bar-Symbionts-5",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-5",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-5",
    "commit": {
      "sha": "c64f40d8258da47602b83f5670c50010471a7de3",
      "url": "https://api.github.com/repos/foo/bar/commits/c64f40d8258da47602b83f5670c50010471a7de3"
    }
  },
  {
    "name": "foo-bar-Symbionts-4",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-4",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-4",
    "commit": {
      "sha": "f96bd08f5bd4286acd637c47eef38ca68c7c8483",
      "url": "https://api.github.com/repos/foo/bar/commits/f96bd08f5bd4286acd637c47eef38ca68c7c8483"
    }
  },
  {
    "name": "foo-bar-Symbionts-3",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-3",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-3",
    "commit": {
      "sha": "bca59f4dbb622a6b1b8673f21b23a20cc1efb472",
      "url": "https://api.github.com/repos/foo/bar/commits/bca59f4dbb622a6b1b8673f21b23a20cc1efb472"
    }
  },
  {
    "name": "foo-bar-Symbionts-2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-2",
    "commit": {
      "sha": "b5775376378dfbc2680017598bb2a64faff39e73",
      "url": "https://api.github.com/repos/foo/bar/commits/b5775376378dfbc2680017598bb2a64faff39e73"
    }
  },
  {
    "name": "foo-bar-Symbionts-1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Symbionts-1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Symbionts-1",
    "commit": {
      "sha": "7de1b288136f2cbb8c5fb458c09e9fbcbe99be34",
      "url": "https://api.github.com/repos/foo/bar/commits/7de1b288136f2cbb8c5fb458c09e9fbcbe99be34"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-44.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-44.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-44.0.0",
    "commit": {
      "sha": "2205d7dc98430506ab13842f55720f2f5319ae3d",
      "url": "https://api.github.com/repos/foo/bar/commits/2205d7dc98430506ab13842f55720f2f5319ae3d"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-43.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-43.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-43.0.0",
    "commit": {
      "sha": "25159d5f23842bed063edb64716e9a73d5b6090b",
      "url": "https://api.github.com/repos/foo/bar/commits/25159d5f23842bed063edb64716e9a73d5b6090b"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-42.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-42.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-42.0.0",
    "commit": {
      "sha": "a9487d96c66aa92d01a44cd62fddd139652b5d20",
      "url": "https://api.github.com/repos/foo/bar/commits/a9487d96c66aa92d01a44cd62fddd139652b5d20"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-38.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-38.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-38.0.0",
    "commit": {
      "sha": "3668c7148546840ed4db1ca3f4bd25dc9391bcb5",
      "url": "https://api.github.com/repos/foo/bar/commits/3668c7148546840ed4db1ca3f4bd25dc9391bcb5"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-37.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-37.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-37.0.1",
    "commit": {
      "sha": "ea13511b7c8f701b002d52c431f3d842bbc7bfe2",
      "url": "https://api.github.com/repos/foo/bar/commits/ea13511b7c8f701b002d52c431f3d842bbc7bfe2"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-37.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-37.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-37.0.0",
    "commit": {
      "sha": "515f3473a63a0d9fa6a941bd3f3941b072048629",
      "url": "https://api.github.com/repos/foo/bar/commits/515f3473a63a0d9fa6a941bd3f3941b072048629"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.7",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.7",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.7",
    "commit": {
      "sha": "d48e3ed53fb61896f67a6cba4c742d4e3fbfae58",
      "url": "https://api.github.com/repos/foo/bar/commits/d48e3ed53fb61896f67a6cba4c742d4e3fbfae58"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.6",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.6",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.6",
    "commit": {
      "sha": "089ef3f318eb566add261a44e5b3354f400d7b1c",
      "url": "https://api.github.com/repos/foo/bar/commits/089ef3f318eb566add261a44e5b3354f400d7b1c"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.5",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.5",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.5",
    "commit": {
      "sha": "ad750eea386ce817643adc2171880552b38f345b",
      "url": "https://api.github.com/repos/foo/bar/commits/ad750eea386ce817643adc2171880552b38f345b"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.4",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.4",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.4",
    "commit": {
      "sha": "2d2ecb93b3e0fbc1cc5fffe66d1f76bcb94585c3",
      "url": "https://api.github.com/repos/foo/bar/commits/2d2ecb93b3e0fbc1cc5fffe66d1f76bcb94585c3"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.3",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.3",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.3",
    "commit": {
      "sha": "8085b6cd144c4bc441644c5626be41e3b0d8e054",
      "url": "https://api.github.com/repos/foo/bar/commits/8085b6cd144c4bc441644c5626be41e3b0d8e054"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.2",
    "commit": {
      "sha": "9760e8f7ff7fb70f25c9cdede81ddb918c325954",
      "url": "https://api.github.com/repos/foo/bar/commits/9760e8f7ff7fb70f25c9cdede81ddb918c325954"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.1",
    "commit": {
      "sha": "3a454014999e83256b39af18d007d82e67ffb09c",
      "url": "https://api.github.com/repos/foo/bar/commits/3a454014999e83256b39af18d007d82e67ffb09c"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-36.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-36.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-36.0.0",
    "commit": {
      "sha": "270b1ab4ab2948c00251ded6cc4927e44f64ef29",
      "url": "https://api.github.com/repos/foo/bar/commits/270b1ab4ab2948c00251ded6cc4927e44f64ef29"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-35.3.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-35.3.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-35.3.0",
    "commit": {
      "sha": "0cb9a3a9cae2727b33a9f73343348fc41fbdf51f",
      "url": "https://api.github.com/repos/foo/bar/commits/0cb9a3a9cae2727b33a9f73343348fc41fbdf51f"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-35.2.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-35.2.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-35.2.0",
    "commit": {
      "sha": "29c3065cb8fb80b82f4c47a5fed696edb5796996",
      "url": "https://api.github.com/repos/foo/bar/commits/29c3065cb8fb80b82f4c47a5fed696edb5796996"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-35.1.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-35.1.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-35.1.0",
    "commit": {
      "sha": "1b90170cb45f28cf9e28d7833c29fb433506efbf",
      "url": "https://api.github.com/repos/foo/bar/commits/1b90170cb45f28cf9e28d7833c29fb433506efbf"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-35.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-35.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-35.0.1",
    "commit": {
      "sha": "3a454014999e83256b39af18d007d82e67ffb09c",
      "url": "https://api.github.com/repos/foo/bar/commits/3a454014999e83256b39af18d007d82e67ffb09c"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-35.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-35.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-35.0.0",
    "commit": {
      "sha": "e5caf71838e793f13130fedc6986d5d8e746e85a",
      "url": "https://api.github.com/repos/foo/bar/commits/e5caf71838e793f13130fedc6986d5d8e746e85a"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-34.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-34.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-34.0.0",
    "commit": {
      "sha": "50d87a0039566526b85db7bee1f8330fa436989b",
      "url": "https://api.github.com/repos/foo/bar/commits/50d87a0039566526b85db7bee1f8330fa436989b"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-33.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-33.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-33.0.0",
    "commit": {
      "sha": "502567da2c24b3995507a6edf64dd88f4d58d162",
      "url": "https://api.github.com/repos/foo/bar/commits/502567da2c24b3995507a6edf64dd88f4d58d162"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-32.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-32.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-32.0.0",
    "commit": {
      "sha": "e0c7bd61b1afa31ea7d5f41780a413e337f7d9f0",
      "url": "https://api.github.com/repos/foo/bar/commits/e0c7bd61b1afa31ea7d5f41780a413e337f7d9f0"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-31.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-31.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-31.0.0",
    "commit": {
      "sha": "03c6d2b43260aed6bb7b8c7e27aa9cdae8df5a32",
      "url": "https://api.github.com/repos/foo/bar/commits/03c6d2b43260aed6bb7b8c7e27aa9cdae8df5a32"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-30.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-30.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-30.0.0",
    "commit": {
      "sha": "afd8d82ae029416cea944135d28b10085fe209d2",
      "url": "https://api.github.com/repos/foo/bar/commits/afd8d82ae029416cea944135d28b10085fe209d2"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-29.1.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-29.1.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-29.1.0",
    "commit": {
      "sha": "a57e068047c4180b1e74e6670e187a616139ae79",
      "url": "https://api.github.com/repos/foo/bar/commits/a57e068047c4180b1e74e6670e187a616139ae79"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-29.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-29.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-29.0.0",
    "commit": {
      "sha": "cd84c4a43c99731a7b51bde203dff276347ffd2d",
      "url": "https://api.github.com/repos/foo/bar/commits/cd84c4a43c99731a7b51bde203dff276347ffd2d"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.5.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.5.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.5.1",
    "commit": {
      "sha": "d548d04dbee49546bddfec72c9078718812b1bf9",
      "url": "https://api.github.com/repos/foo/bar/commits/d548d04dbee49546bddfec72c9078718812b1bf9"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.5.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.5.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.5.0",
    "commit": {
      "sha": "f408c511f18b36dcc44f909367ff88886824fa45",
      "url": "https://api.github.com/repos/foo/bar/commits/f408c511f18b36dcc44f909367ff88886824fa45"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.4.3",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.4.3",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.4.3",
    "commit": {
      "sha": "402e7b914c15890a9f5ea4617da673b906211260",
      "url": "https://api.github.com/repos/foo/bar/commits/402e7b914c15890a9f5ea4617da673b906211260"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.4.2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.4.2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.4.2",
    "commit": {
      "sha": "98bb0264db9c43257ca2d114f16daf74185fd547",
      "url": "https://api.github.com/repos/foo/bar/commits/98bb0264db9c43257ca2d114f16daf74185fd547"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.4.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.4.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.4.1",
    "commit": {
      "sha": "4a59bcdb0f9b40854efb5385ba140e0bb8b82b1e",
      "url": "https://api.github.com/repos/foo/bar/commits/4a59bcdb0f9b40854efb5385ba140e0bb8b82b1e"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.4.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.4.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.4.0",
    "commit": {
      "sha": "63edb33ecd6c6232defa06aae45ee2d43fab0ccf",
      "url": "https://api.github.com/repos/foo/bar/commits/63edb33ecd6c6232defa06aae45ee2d43fab0ccf"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.3.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.3.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.3.1",
    "commit": {
      "sha": "2db48a5d2fe52435394e1d418b3e7b67111a7904",
      "url": "https://api.github.com/repos/foo/bar/commits/2db48a5d2fe52435394e1d418b3e7b67111a7904"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.3.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.3.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.3.0",
    "commit": {
      "sha": "afd0cbf85c0f2c0f1debe717bb8aa21d78ed860a",
      "url": "https://api.github.com/repos/foo/bar/commits/afd0cbf85c0f2c0f1debe717bb8aa21d78ed860a"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.2.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.2.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.2.1",
    "commit": {
      "sha": "dcc3cf2727813eb51a98b94c699d6c53c3be9ea0",
      "url": "https://api.github.com/repos/foo/bar/commits/dcc3cf2727813eb51a98b94c699d6c53c3be9ea0"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.2.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.2.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.2.0",
    "commit": {
      "sha": "fc955c9b287ebfbce5824e9b7712ab11b9d0785c",
      "url": "https://api.github.com/repos/foo/bar/commits/fc955c9b287ebfbce5824e9b7712ab11b9d0785c"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.1.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.1.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.1.0",
    "commit": {
      "sha": "539c345687e200aec22353af6b73391dd245ed09",
      "url": "https://api.github.com/repos/foo/bar/commits/539c345687e200aec22353af6b73391dd245ed09"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.0.2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.0.2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.0.2",
    "commit": {
      "sha": "28dd0e8640f2fd597ede6a04be7dcbc2e16df558",
      "url": "https://api.github.com/repos/foo/bar/commits/28dd0e8640f2fd597ede6a04be7dcbc2e16df558"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.0.1",
    "commit": {
      "sha": "63d38b3cefca89443a7f1477a204e83fae41c086",
      "url": "https://api.github.com/repos/foo/bar/commits/63d38b3cefca89443a7f1477a204e83fae41c086"
    }
  },
  {
    "name": "foo-bar-Sprint-Release-28.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-Release-28.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-Release-28.0.0",
    "commit": {
      "sha": "473ead490ae19f3fba9c09f2c5b56661b0460eae",
      "url": "https://api.github.com/repos/foo/bar/commits/473ead490ae19f3fba9c09f2c5b56661b0460eae"
    }
  },
  {
    "name": "foo-bar-Sprint-29.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-29.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-29.0.0",
    "commit": {
      "sha": "25afeea35d9a8fb5eabd7b6e6cceab89be196e84",
      "url": "https://api.github.com/repos/foo/bar/commits/25afeea35d9a8fb5eabd7b6e6cceab89be196e84"
    }
  },
  {
    "name": "foo-bar-Sprint-28.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-28.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-28.0.0",
    "commit": {
      "sha": "c1360c0dbf6dbf319007d6a189b8b836b2771b30",
      "url": "https://api.github.com/repos/foo/bar/commits/c1360c0dbf6dbf319007d6a189b8b836b2771b30"
    }
  },
  {
    "name": "foo-bar-Sprint-28.0.0-beta",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-28.0.0-beta",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-28.0.0-beta",
    "commit": {
      "sha": "0e1486a849694fc454382b77b1bf0ff1d0b946bf",
      "url": "https://api.github.com/repos/foo/bar/commits/0e1486a849694fc454382b77b1bf0ff1d0b946bf"
    }
  },
  {
    "name": "foo-bar-Sprint-28-base",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-28-base",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-28-base",
    "commit": {
      "sha": "9e774ec33584fb9c23810e9e337fa036978dbe1d",
      "url": "https://api.github.com/repos/foo/bar/commits/9e774ec33584fb9c23810e9e337fa036978dbe1d"
    }
  },
  {
    "name": "foo-bar-Sprint-27.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-27.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-27.0.0",
    "commit": {
      "sha": "9e774ec33584fb9c23810e9e337fa036978dbe1d",
      "url": "https://api.github.com/repos/foo/bar/commits/9e774ec33584fb9c23810e9e337fa036978dbe1d"
    }
  },
  {
    "name": "foo-bar-Sprint-26.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-26.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-26.0.0",
    "commit": {
      "sha": "2c476b032a5e2597147ab257675726484b7a641b",
      "url": "https://api.github.com/repos/foo/bar/commits/2c476b032a5e2597147ab257675726484b7a641b"
    }
  },
  {
    "name": "foo-bar-Sprint-24.0.4",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-24.0.4",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-24.0.4",
    "commit": {
      "sha": "2121bab7c955214f2671c4d66ae1113551b85ff3",
      "url": "https://api.github.com/repos/foo/bar/commits/2121bab7c955214f2671c4d66ae1113551b85ff3"
    }
  },
  {
    "name": "foo-bar-Sprint-24.0.3",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-24.0.3",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-24.0.3",
    "commit": {
      "sha": "d520038d0e83c75ba2f11ce2c35957853ea1b314",
      "url": "https://api.github.com/repos/foo/bar/commits/d520038d0e83c75ba2f11ce2c35957853ea1b314"
    }
  },
  {
    "name": "foo-bar-Sprint-24.0.2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-24.0.2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-24.0.2",
    "commit": {
      "sha": "f7c0b52571db224b07f2f13ab5a11cee9cc7cc08",
      "url": "https://api.github.com/repos/foo/bar/commits/f7c0b52571db224b07f2f13ab5a11cee9cc7cc08"
    }
  },
  {
    "name": "foo-bar-Sprint-24.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-24.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-24.0.0",
    "commit": {
      "sha": "f2f309977537e2c23556627b32727d3d45aa6728",
      "url": "https://api.github.com/repos/foo/bar/commits/f2f309977537e2c23556627b32727d3d45aa6728"
    }
  },
  {
    "name": "foo-bar-Sprint-23.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-23.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-23.0.0",
    "commit": {
      "sha": "f033bcf0c06695a72c5e9ca1b6d7a8e650379799",
      "url": "https://api.github.com/repos/foo/bar/commits/f033bcf0c06695a72c5e9ca1b6d7a8e650379799"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.8",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.8",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.8",
    "commit": {
      "sha": "6099db8b528536a8394b2870e8c28ffa44ad115a",
      "url": "https://api.github.com/repos/foo/bar/commits/6099db8b528536a8394b2870e8c28ffa44ad115a"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.7",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.7",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.7",
    "commit": {
      "sha": "00b03e8cb4e02270744a97e0d1dd3447aca2d4e2",
      "url": "https://api.github.com/repos/foo/bar/commits/00b03e8cb4e02270744a97e0d1dd3447aca2d4e2"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.6",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.6",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.6",
    "commit": {
      "sha": "c55a33c9b891849dbedfca235d218823664f6628",
      "url": "https://api.github.com/repos/foo/bar/commits/c55a33c9b891849dbedfca235d218823664f6628"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.5",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.5",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.5",
    "commit": {
      "sha": "8c4557c89f77547df5575650c575954662b054f4",
      "url": "https://api.github.com/repos/foo/bar/commits/8c4557c89f77547df5575650c575954662b054f4"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.4",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.4",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.4",
    "commit": {
      "sha": "8c4557c89f77547df5575650c575954662b054f4",
      "url": "https://api.github.com/repos/foo/bar/commits/8c4557c89f77547df5575650c575954662b054f4"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.3",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.3",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.3",
    "commit": {
      "sha": "1dd23e5f1aa24a38cc1bfe6ecc97822d785fad14",
      "url": "https://api.github.com/repos/foo/bar/commits/1dd23e5f1aa24a38cc1bfe6ecc97822d785fad14"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.2",
    "commit": {
      "sha": "1dd23e5f1aa24a38cc1bfe6ecc97822d785fad14",
      "url": "https://api.github.com/repos/foo/bar/commits/1dd23e5f1aa24a38cc1bfe6ecc97822d785fad14"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.1",
    "commit": {
      "sha": "1dd23e5f1aa24a38cc1bfe6ecc97822d785fad14",
      "url": "https://api.github.com/repos/foo/bar/commits/1dd23e5f1aa24a38cc1bfe6ecc97822d785fad14"
    }
  },
  {
    "name": "foo-bar-Sprint-21.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-21.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-21.0.0",
    "commit": {
      "sha": "9c841e33207130e36f7ce03a93c7de98b2ed5ae9",
      "url": "https://api.github.com/repos/foo/bar/commits/9c841e33207130e36f7ce03a93c7de98b2ed5ae9"
    }
  },
  {
    "name": "foo-bar-Sprint-20.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-20.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-20.0.1",
    "commit": {
      "sha": "45e03e9c9ef5eb1de9325a3801f49c48a47718f7",
      "url": "https://api.github.com/repos/foo/bar/commits/45e03e9c9ef5eb1de9325a3801f49c48a47718f7"
    }
  },
  {
    "name": "foo-bar-Sprint-20.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-20.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-20.0.0",
    "commit": {
      "sha": "38e723df8267e2a4e9c9c3d8571cfed4b783c7ec",
      "url": "https://api.github.com/repos/foo/bar/commits/38e723df8267e2a4e9c9c3d8571cfed4b783c7ec"
    }
  },
  {
    "name": "foo-bar-Sprint-19.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-19.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-19.0.1",
    "commit": {
      "sha": "90e87f56aee589ca4c9bc85532bb9330ff0f9bee",
      "url": "https://api.github.com/repos/foo/bar/commits/90e87f56aee589ca4c9bc85532bb9330ff0f9bee"
    }
  },
  {
    "name": "foo-bar-Sprint-19.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-19.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-19.0.0",
    "commit": {
      "sha": "ec3937ba699912d88ac954719c5ab8916257178f",
      "url": "https://api.github.com/repos/foo/bar/commits/ec3937ba699912d88ac954719c5ab8916257178f"
    }
  },
  {
    "name": "foo-bar-Sprint-18.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-18.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-18.0.1",
    "commit": {
      "sha": "72bbf49920473aed11199c05c0e92eb097eed2dc",
      "url": "https://api.github.com/repos/foo/bar/commits/72bbf49920473aed11199c05c0e92eb097eed2dc"
    }
  },
  {
    "name": "foo-bar-Sprint-18.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-18.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-18.0.0",
    "commit": {
      "sha": "a7cbeee241215cbde367e53f2fe3372a337c60b9",
      "url": "https://api.github.com/repos/foo/bar/commits/a7cbeee241215cbde367e53f2fe3372a337c60b9"
    }
  },
  {
    "name": "foo-bar-Sprint-17.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-17.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-17.0.0",
    "commit": {
      "sha": "7896dde32ee24111be5717377ab69febf82af364",
      "url": "https://api.github.com/repos/foo/bar/commits/7896dde32ee24111be5717377ab69febf82af364"
    }
  },
  {
    "name": "foo-bar-Sprint-16.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-16.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-16.0.0",
    "commit": {
      "sha": "9ac1c8bc4f1314ed076e23e642d7285688c8aec4",
      "url": "https://api.github.com/repos/foo/bar/commits/9ac1c8bc4f1314ed076e23e642d7285688c8aec4"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.10",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.10",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.10",
    "commit": {
      "sha": "c5dcc51c46b48489d7bacb8d6381632ec03e184e",
      "url": "https://api.github.com/repos/foo/bar/commits/c5dcc51c46b48489d7bacb8d6381632ec03e184e"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.9",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.9",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.9",
    "commit": {
      "sha": "228e17f0e00815631c75c9721cc687b04db8bdb1",
      "url": "https://api.github.com/repos/foo/bar/commits/228e17f0e00815631c75c9721cc687b04db8bdb1"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.8",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.8",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.8",
    "commit": {
      "sha": "46777ee0b6d87fde267bf77aa96b50f388524cde",
      "url": "https://api.github.com/repos/foo/bar/commits/46777ee0b6d87fde267bf77aa96b50f388524cde"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.7",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.7",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.7",
    "commit": {
      "sha": "4b22d344913a1a509c372155af2f708a4b71c62f",
      "url": "https://api.github.com/repos/foo/bar/commits/4b22d344913a1a509c372155af2f708a4b71c62f"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.6",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.6",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.6",
    "commit": {
      "sha": "2a4b3fdb772f8be58a337a54d689f98462a54a30",
      "url": "https://api.github.com/repos/foo/bar/commits/2a4b3fdb772f8be58a337a54d689f98462a54a30"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.5",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.5",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.5",
    "commit": {
      "sha": "997419a25b4548e43caf66f95d2809611ade0cb4",
      "url": "https://api.github.com/repos/foo/bar/commits/997419a25b4548e43caf66f95d2809611ade0cb4"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.4",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.4",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.4",
    "commit": {
      "sha": "a6a67f8f8adbe9f11e2b6a60d9fd3a713abb046b",
      "url": "https://api.github.com/repos/foo/bar/commits/a6a67f8f8adbe9f11e2b6a60d9fd3a713abb046b"
    }
  },
  {
    "name": "foo-bar-Sprint-15.0.3",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.3",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.3",
    "commit": {
      "sha": "c61d1695cd9c9ee8cc890159b16f28028adbc683",
      "url": "https://api.github.com/repos/foo/bar/commits/c61d1695cd9c9ee8cc890159b16f28028adbc683"
    }
  },
  {
    "name": "15.0.2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.2",
    "commit": {
      "sha": "5b6837a4755d2375e56ed019d79d53913eb00368",
      "url": "https://api.github.com/repos/foo/bar/commits/5b6837a4755d2375e56ed019d79d53913eb00368"
    }
  },
  {
    "name": "v15.0.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.1",
    "commit": {
      "sha": "0cb8b7790295c9095403da312a38d6a7266bf6bc",
      "url": "https://api.github.com/repos/foo/bar/commits/0cb8b7790295c9095403da312a38d6a7266bf6bc"
    }
  },
  {
    "name": "15.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-15.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-15.0.0",
    "commit": {
      "sha": "c50e8f6b7f0490fb63105b85c507fd19d7483263",
      "url": "https://api.github.com/repos/foo/bar/commits/c50e8f6b7f0490fb63105b85c507fd19d7483263"
    }
  },
  {
    "name": "14.1.18",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-14.1.18",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-14.1.18",
    "commit": {
      "sha": "3535e75de62c14059208dd0e50d07a2e4032dc2f",
      "url": "https://api.github.com/repos/foo/bar/commits/3535e75de62c14059208dd0e50d07a2e4032dc2f"
    }
  },
  {
    "name": "v14.1.17",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-14.1.17",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-14.1.17",
    "commit": {
      "sha": "1249d2cfbbed951db47bba45d0b27805017c369d",
      "url": "https://api.github.com/repos/foo/bar/commits/1249d2cfbbed951db47bba45d0b27805017c369d"
    }
  },
  {
    "name": "v14.1.16",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/foo-bar-Sprint-14.1.16",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/foo-bar-Sprint-14.1.16",
    "commit": {
      "sha": "23ca259078620d135326dd360ba79e10da4c12f5",
      "url": "https://api.github.com/repos/foo/bar/commits/23ca259078620d135326dd360ba79e10da4c12f5"
    }
  }
]
`

const apiTagsPage1Link = `
<https://api.github.com/repositories/12345678/tags?per_page=100&page=2>; rel="next", <https://api.github.com/repositories/12345678/tags?per_page=100&page=2>; rel="last"
`

const apiTagsPage2 = `
[
  {
    "name": "DTC-1338-2",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/DTC-1338-2",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/DTC-1338-2",
    "commit": {
      "sha": "71bb00261668280b11e8833a03565c03e91b4be8",
      "url": "https://api.github.com/repos/foo/bar/commits/71bb00261668280b11e8833a03565c03e91b4be8"
    }
  },
  {
    "name": "DTC-1338-1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/DTC-1338-1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/DTC-1338-1",
    "commit": {
      "sha": "bbe1ff2f47c62f24d2f6a435846b7261303ba040",
      "url": "https://api.github.com/repos/foo/bar/commits/bbe1ff2f47c62f24d2f6a435846b7261303ba040"
    }
  },
  {
    "name": "46.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/46.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/46.0.0",
    "commit": {
      "sha": "1efe511365e7ebdf7e79026121254fbbe4f79945",
      "url": "https://api.github.com/repos/foo/bar/commits/1efe511365e7ebdf7e79026121254fbbe4f79945"
    }
  },
  {
    "name": "1.1.1",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/1.1.1",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/1.1.1",
    "commit": {
      "sha": "1efe511365e7ebdf7e79026121254fbbe4f79945",
      "url": "https://api.github.com/repos/foo/bar/commits/1efe511365e7ebdf7e79026121254fbbe4f79945"
    }
  },
  {
    "name": "1.0.0",
    "zipball_url": "https://api.github.com/repos/foo/bar/zipball/1.0.0",
    "tarball_url": "https://api.github.com/repos/foo/bar/tarball/1.0.0",
    "commit": {
      "sha": "1efe511365e7ebdf7e79026121254fbbe4f79945",
      "url": "https://api.github.com/repos/foo/bar/commits/1efe511365e7ebdf7e79026121254fbbe4f79945"
    }
  }
]
`

const apiTagsPage2Link = `
<https://api.github.com/repositories/12345678/tags?per_page=100&page=1>; rel="prev", <https://api.github.com/repositories/12345678/tags?per_page=100&page=1>; rel="first"
`

var apiTagsExpected = []string{
	"tag-jin-tries-arbitrary-tag-01",
	"tag-non-semantic-version-1.0.05.6",
	"tag-jin-prefix-0.0.1",
	"tag-fix-hcc-iframe-for-ssl",
	"tag-fix-cid-generation",
	"foo-bar-Symbionts",
	"foo-bar-Symbionts-12",
	"foo-bar-Symbionts-11",
	"foo-bar-Symbionts-10",
	"foo-bar-Symbionts-9",
	"foo-bar-Symbionts-8",
	"foo-bar-Symbionts-7",
	"foo-bar-Symbionts-6",
	"foo-bar-Symbionts-5",
	"foo-bar-Symbionts-4",
	"foo-bar-Symbionts-3",
	"foo-bar-Symbionts-2",
	"foo-bar-Symbionts-1",
	"foo-bar-Sprint-Release-44.0.0",
	"foo-bar-Sprint-Release-43.0.0",
	"foo-bar-Sprint-Release-42.0.0",
	"foo-bar-Sprint-Release-38.0.0",
	"foo-bar-Sprint-Release-37.0.1",
	"foo-bar-Sprint-Release-37.0.0",
	"foo-bar-Sprint-Release-36.0.7",
	"foo-bar-Sprint-Release-36.0.6",
	"foo-bar-Sprint-Release-36.0.5",
	"foo-bar-Sprint-Release-36.0.4",
	"foo-bar-Sprint-Release-36.0.3",
	"foo-bar-Sprint-Release-36.0.2",
	"foo-bar-Sprint-Release-36.0.1",
	"foo-bar-Sprint-Release-36.0.0",
	"foo-bar-Sprint-Release-35.3.0",
	"foo-bar-Sprint-Release-35.2.0",
	"foo-bar-Sprint-Release-35.1.0",
	"foo-bar-Sprint-Release-35.0.1",
	"foo-bar-Sprint-Release-35.0.0",
	"foo-bar-Sprint-Release-34.0.0",
	"foo-bar-Sprint-Release-33.0.0",
	"foo-bar-Sprint-Release-32.0.0",
	"foo-bar-Sprint-Release-31.0.0",
	"foo-bar-Sprint-Release-30.0.0",
	"foo-bar-Sprint-Release-29.1.0",
	"foo-bar-Sprint-Release-29.0.0",
	"foo-bar-Sprint-Release-28.5.1",
	"foo-bar-Sprint-Release-28.5.0",
	"foo-bar-Sprint-Release-28.4.3",
	"foo-bar-Sprint-Release-28.4.2",
	"foo-bar-Sprint-Release-28.4.1",
	"foo-bar-Sprint-Release-28.4.0",
	"foo-bar-Sprint-Release-28.3.1",
	"foo-bar-Sprint-Release-28.3.0",
	"foo-bar-Sprint-Release-28.2.1",
	"foo-bar-Sprint-Release-28.2.0",
	"foo-bar-Sprint-Release-28.1.0",
	"foo-bar-Sprint-Release-28.0.2",
	"foo-bar-Sprint-Release-28.0.1",
	"foo-bar-Sprint-Release-28.0.0",
	"foo-bar-Sprint-29.0.0",
	"foo-bar-Sprint-28.0.0",
	"foo-bar-Sprint-28.0.0-beta",
	"foo-bar-Sprint-28-base",
	"foo-bar-Sprint-27.0.0",
	"foo-bar-Sprint-26.0.0",
	"foo-bar-Sprint-24.0.4",
	"foo-bar-Sprint-24.0.3",
	"foo-bar-Sprint-24.0.2",
	"foo-bar-Sprint-24.0.0",
	"foo-bar-Sprint-23.0.0",
	"foo-bar-Sprint-21.0.8",
	"foo-bar-Sprint-21.0.7",
	"foo-bar-Sprint-21.0.6",
	"foo-bar-Sprint-21.0.5",
	"foo-bar-Sprint-21.0.4",
	"foo-bar-Sprint-21.0.3",
	"foo-bar-Sprint-21.0.2",
	"foo-bar-Sprint-21.0.1",
	"foo-bar-Sprint-21.0.0",
	"foo-bar-Sprint-20.0.1",
	"foo-bar-Sprint-20.0.0",
	"foo-bar-Sprint-19.0.1",
	"foo-bar-Sprint-19.0.0",
	"foo-bar-Sprint-18.0.1",
	"foo-bar-Sprint-18.0.0",
	"foo-bar-Sprint-17.0.0",
	"foo-bar-Sprint-16.0.0",
	"foo-bar-Sprint-15.0.10",
	"foo-bar-Sprint-15.0.9",
	"foo-bar-Sprint-15.0.8",
	"foo-bar-Sprint-15.0.7",
	"foo-bar-Sprint-15.0.6",
	"foo-bar-Sprint-15.0.5",
	"foo-bar-Sprint-15.0.4",
	"foo-bar-Sprint-15.0.3",
	"15.0.2",
	"v15.0.1",
	"15.0.0",
	"14.1.18",
	"v14.1.17",
	"v14.1.16",
	"DTC-1338-2",
	"DTC-1338-1",
	"46.0.0",
	"1.1.1",
	"1.0.0",
}

const relsPage1 = `
[
  {
    "url": "https://api.github.com/repos/sna/fu/releases/9876543",
    "id": 9876543,
    "tag_name": "9.8.7",
    "name": "prerelease",
    "prerelease": true,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/9854782",
        "id": 9854782,
        "name": "foo.tgz"
      },
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/9854783",
        "id": 9854783,
        "name": "bar.tgz"
      }
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/8765432",
    "id": 8765432,
    "tag_name": "bad8.7.6.5",
    "name": "bad tag",
    "prerelease": false,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/8754782",
        "id": 8754782,
        "name": "foo.tgz"
      },
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/8754783",
        "id": 8754783,
        "name": "bar.tgz"
      }
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/7654321",
    "id": 7654321,
    "tag_name": "7.6.5",
    "name": "all good",
    "prerelease": false,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/7654782",
        "id": 7654782,
        "name": "foo.tgz"
      },
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/7654783",
        "id": 7654783,
        "name": "bar.tgz"
      }
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/6556789",
    "id": 6556789,
    "tag_name": "6.5.5",
    "name": "missing asset",
    "prerelease": false,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/6554782",
        "id": 6554782,
        "name": "foo.tgz"
      }
    ]
  }
]
`

const relsPage2 = `
[
  {
    "url": "https://api.github.com/repos/sna/fu/releases/1234567",
    "id": 1234567,
    "tag_name": "v1.2.3",
    "name": "prerelease",
    "prerelease": true,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/5354782",
        "id": 5354782,
        "name": "foo.tgz"
      },
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/5354783",
        "id": 5354783,
        "name": "bar.tgz"
      }
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/2345678",
    "id": 2345678,
    "tag_name": "bad2.3.4.0",
    "name": "bad tag",
    "prerelease": false,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/6354782",
        "id": 6354782,
        "name": "foo.tgz"
      },
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/6354783",
        "id": 6354783,
        "name": "bar.tgz"
      }
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/3456789",
    "id": 3456789,
    "tag_name": "3.4.5",
    "name": "all good",
    "prerelease": false,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/7354782",
        "id": 7354782,
        "name": "foo.tgz"
      },
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/7354783",
        "id": 7354783,
        "name": "bar.tgz"
      }
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/4556789",
    "id": 4556789,
    "tag_name": "4.5.6",
    "name": "missing asset",
    "prerelease": false,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/8354782",
        "id": 8354782,
        "name": "foo.tgz"
      }
    ]
  }
]
`

const relsPage1Link = `
<https://api.github.com/repositories/12345678/releases?per_page=100&page=2>; rel="next", <https://api.github.com/repositories/12345678/releases?per_page=100&page=2>; rel="last"
`

const relsPage2Link = `
<https://api.github.com/repositories/12345678/releases?per_page=100&page=1>; rel="prev", <https://api.github.com/repositories/12345678/releases?per_page=100&page=1>; rel="first"
`

const noValidRelsPage1 = `
[
  {
    "url": "https://api.github.com/repos/sna/fu/releases/9876543",
    "id": 9876543,
    "tag_name": "9.8.7",
    "name": "prerelease, no assets",
    "prerelease": true,
    "assets": [
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/8765432",
    "id": 8765432,
    "tag_name": "bad8.7.6.5",
    "name": "another prerelease, no assets",
    "prerelease": true,
    "assets": [
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/7654321",
    "id": 7654321,
    "tag_name": "7.6.5",
    "name": "yet another prerelease, no assets",
    "prerelease": false,
    "assets": [
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/6556789",
    "id": 6556789,
    "tag_name": "6.5.5",
    "name": "seriously, not a single attached release? And you still haven't fixed this?",
    "prerelease": false,
    "assets": [
    ]
  }
]
`
const noValidRelsPage2 = `
[
  {
    "url": "https://api.github.com/repos/sna/fu/releases/1234567",
    "id": 1234567,
    "tag_name": "v1.2.3",
    "name": "C'mon. Just one little release attachment! You can do it! I believe IN YOU!!!",
    "prerelease": true,
    "assets": [
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/2345678",
    "id": 2345678,
    "tag_name": "bad2.3.4.0",
    "name": "Ah go on! Do it for the children!",
    "prerelease": true,
    "assets": [
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/3456789",
    "id": 3456789,
    "tag_name": "3.4.5",
    "name": "Are you freakin' kidding me? evil-virus.tgz!?",
    "prerelease": true,
    "assets": [
      {
        "url": "https://api.github.com/repos/sna/fu/releases/assets/7354782",
        "id": 7354782,
        "name": "evil-virus.tgz"
      },
    ]
  },
  {
    "url": "https://api.github.com/repos/sna/fu/releases/4556789",
    "id": 4556789,
    "tag_name": "4.5.6",
    "name": "Screw you guys, I'm going home.",
    "prerelease": true,
    "assets": [
    ]
  }
]
`
const noValidRelsPage1Link = `
<https://api.github.com/repositories/12345678/releases?per_page=100&page=2>; rel="next", <https://api.github.com/repositories/12345678/releases?per_page=100&page=2>; rel="last"
`

const noValidRelsPage2Link = `
<https://api.github.com/repositories/12345678/releases?per_page=100&page=1>; rel="prev", <https://api.github.com/repositories/12345678/releases?per_page=100&page=1>; rel="first"
`
