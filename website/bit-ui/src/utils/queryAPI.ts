import axios from "axios";
import dayjs from "dayjs";

const BIT_QUERY_URL = 'http://localhost:8085';
const dateFormat = 'YYYY-MMMM-DD HH:mm:s';

export const getBitStatus = async (userGroup: string, startTime: Date, endTime: Date, filter: string) => {
    const start = dayjs(startTime).format(dateFormat);
    const end = dayjs(endTime).format(dateFormat);
    const res = await axios.get(BIT_QUERY_URL + '/status',
        {params: {user_group: userGroup, start: start , end: end, filter: filter}})
    return res.data;
}

export const getReports = async (filter: string, startTime?: Date, endTime?: Date) => {
    const start = dayjs(startTime).format(dateFormat);
    const end = dayjs(endTime).format(dateFormat);
    const res = await axios.get(BIT_QUERY_URL + '/reports',
        {params: {filter: filter, start: start, end: end}})
    return res.data;
}
