# equips
A way of adding lists of tokens onchain, processed via an indexer. It's like erc6551, but the tokens don't have ownership.
![](https://media.discordapp.net/attachments/1091126083895693312/1123754325798301807/morganmoskalyk_engineered_inventory_systems_with_neon_assembled_99f78f50-321c-4598-a639-b3500beaab95.png?width=800&height=800)
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
