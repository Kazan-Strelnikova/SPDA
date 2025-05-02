export function getDateString(date: Date) {
    const months = [
      'января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
      'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря'
    ];
    const weekdays = [
      'воскресенье', 'понедельник', 'вторник', 'среда',
      'четверг', 'пятница', 'суббота'
    ];
  
    const now = new Date();
    const isToday =
      date.getDate() === now.getDate() &&
      date.getMonth() === now.getMonth() &&
      date.getFullYear() === now.getFullYear();
  
    const day = date.getDate();
    const monthName = months[date.getMonth()];
    const weekdayName = isToday ? 'сегодня' : weekdays[date.getDay()];
  
    return `${day} ${monthName}, ${weekdayName}`;
  }
  