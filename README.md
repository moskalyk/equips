# equips
A way of adding lists of tokens onchain, processed via an indexer. It's like erc6551, but the tokens don't have ownership.

## use cases (similiar to erc6551)
- fashionoutfits
- inventory
- curation
- kabbalah

```js

import equipsIndexer from '/sdk/'

const result1 = await equipsIndexer.all()
const result2 = await equipsIndexer.asOwner("<account_address>")

result1.map((token) => {
	console.log(token.Owner)
	console.log(token.TokenAddress)
	console.log(token.TokenID)
	console.log(token.Index)
	console.log(token.Salt)
})

```
