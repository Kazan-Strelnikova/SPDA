import { Category } from "../types";
import ConcertIcon from "../assets/concert-icon.svg";
import ConferenceIcon from "../assets/conference-icon.svg";
import EducationIcon from "../assets/education-icon.svg";
import ExhibitionIcon from "../assets/exhibition-icon.svg";
import MeetupIcon from "../assets/meetup-icon.svg";
import CompetitionIcon from "../assets/olympiad-icon.svg";
import OtherIcon from "../assets/other-icon.svg";
import PartyIcon from "../assets/party-icon.svg";
import SportIcon from "../assets/sport-icon.svg";


export const getCategoryIcon = (category : Category) => {
    switch (category) {
        case "Conference": return <img src={ConferenceIcon} alt="ConferenceIcon" />;
        case "Meetup": return <img src={MeetupIcon} alt="MeetupIcon" />;
        case "Concert": return <img src={ConcertIcon} alt="ConcertIcon" />;
        case "Exhibition": return <img src={ExhibitionIcon} alt="ExhibitionIcon" />;
        case "Party": return <img src={PartyIcon} alt="PartyIcon" />;
        case "Sport": return <img src={SportIcon} alt="SportIcon" />;
        case "Education": return <img src={EducationIcon} alt="EducationIcon" />;
        case "Competition": return <img src={CompetitionIcon} alt="CompetitionIcon" />;
        case "Other": return <img src={OtherIcon} alt="OtherIcon" />;
    }
}