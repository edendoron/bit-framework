import axios from "axios";
import dayjs from "dayjs";

const BIT_QUERY_URL = 'http://localhost:8085';

export const getBitStatus = async (userGroup: string, startTime: Date, endTime: Date, filter: string) => {
    const res = await axios.get(BIT_QUERY_URL + '/status',
        {params: {user_group: userGroup, start:  dayjs(startTime).format("YYYY-MMMM-DD HH:mm:s"), end: dayjs(endTime).format("YYYY-MMMM-DD HH:mm:s"), filter: filter}})
    return res.data;
}

export const getReports = async (filter: string, startTime?: Date, endTime?: Date) => {
    const res = await axios.get(BIT_QUERY_URL + '/reports',
        {params: {filter: filter, start: dayjs(startTime).format("YYYY-MMMM-DD HH:mm:s"), end: dayjs(endTime).format("YYYY-MMMM-DD HH:mm:s")}})
    return res.data;
}
