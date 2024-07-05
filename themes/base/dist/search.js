var indexData = null

function loadIndex() {
  return fetch("/index.json")
    .then(response => response.json())
    .then(data => {
      indexData = data
    })
}

function tagged(tag) {
  const t = indexData.tm[tag]
  if (t === undefined) {
    return []
  }
  var hits = []
  for (const id of t) {
    hits.push(indexData.docs[id])
  }
  return hits
}

function search(query) {

  query = query.toLowerCase()

  console.log("index is", indexData)
  // split query by space
  var parts = query.split(" ")
  console.log("parts are ", parts)
  console.log("query is ", query)


  var hitMap = {}
  for (var p of parts) {
    const pmap = indexData.kw[p]
    if (pmap === undefined) {
      continue // no hits
    }
    for (const pid of Object.keys(pmap)) {
      if (hitMap[pid] === undefined) {
        hitMap[pid] = 0
      }
      hitMap[pid] += pmap[pid]
    }
  }

  console.log("hitMap is ", hitMap)


  var sortedHitIds = Object.keys(hitMap).sort(function (a, b) {
    return hitMap[b] - hitMap[a]
  })
  console.log("sortedHitIds is ", sortedHitIds)

  var hits = []
  for (var i = 0; i < sortedHitIds.length; i++) {
    var id = sortedHitIds[i]
    var hit = indexData.docs[id]
    hits.push(hit)
  }
  console.log("hits are ", hits)
  return hits
}

// loadIndex().then(() => {
//   console.log("search index loaded")
// })