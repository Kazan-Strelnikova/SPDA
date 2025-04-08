export const getCategoryFromNumber = (categoryNum : number ) => {
    switch (categoryNum) {
        case 0: return "Conference"
        case 1: return "Meetup"
        case 2: return "Concert"
        case 3: return "Exhibition"
        case 4: return "Party"
        case 5: return "Sport"
        case 6: return "Education"
        case 7: return "Competition"
        default: return "Other"
    }
}