function rechercher() {
    let input = document.getElementById("searchInput").value.toLowerCase()
    let resultats = document.getElementById("resultats")

    resultats.innerHTML = ""

    tableauAvecTousLesArtistes.forEach(artist => {
        let found = false


        if (artist.name.toLowerCase().includes(input)) {
            found = true
        }

   
        Object.entries(artist.relations).forEach(([lieu, dates]) => {
            if (lieu.toLowerCase().includes(input)) {
                found = true
            }
            dates.forEach(date => {
                if (date.includes(input)) {
                    found = true
                }
            })
        })

        if (found) {
            let div = document.createElement("div")
            div.innerHTML = `<h3>${artist.name}</h3>`
            resultats.appendChild(div)
        }
    })
}
