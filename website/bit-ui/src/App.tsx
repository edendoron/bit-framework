import React, {useEffect, useState} from 'react';
import Axios from 'axios';
import './App.css';
import {Box, Card, CardContent, CardHeader, createStyles, Grid, makeStyles} from '@material-ui/core'
import {Selector} from "./components/Selector";
import {DatePicker} from "./components/DatePicker";

const STORAGE_DATA_READ_URL = 'http://localhost:8082/data/read';
const queryTypes = ['Reports', 'BIT Status', 'Config Failures'];
const userGroups = ['group1', 'group2', 'group3', 'group4', 'groupRafael', 'TemperatureCelsius group', 'group general', 'groupField'];
const filterOptions = ['time', 'tag', 'field'];
// const currentDate = new Date();
// const fullDate = currentDate.getFullYear() + currentDate.getMonth() +

const useStyles = makeStyles(() =>
    createStyles({
        dateGrid: {
            marginTop: 20,
        },
    },
));

export const App = () => {
    const classes = useStyles();

    const [queryType, setQueryType] = useState('');
    const [userGroup, setUserGroup] = useState('');
    const [filter, setFilter] = useState('');
    const [isDisabled, setDisabled] = useState(true);
    const [startTime, setStartTime] = useState(new Date());
    const [endTime, setEndTime] = useState(new Date());
    const [data, setData] = useState();

    const changeQueryType = (event: React.ChangeEvent<{ value: unknown }>) => {
        setQueryType(event.target.value as string);
        setFilter('');
    }

    const changeUserGroup = (event: React.ChangeEvent<{ value: unknown }>) => {
        setUserGroup(event.target.value as string);
        setFilter('');
        setQueryType('');
        setDisabled(false);
    }

    const changeFilter = (event: React.ChangeEvent<{ value: unknown }>) => {
        setFilter((event.target.value as string));
    }

    const changeStartTime = (event: React.ChangeEvent<{ value: unknown }>) => {
        setStartTime(event.target.value as Date)
    }

    const changeEndTime = (event: React.ChangeEvent<{ value: unknown }>) => {
        setEndTime(event.target.value as Date)
    }

    useEffect(() => {
        Axios.get(STORAGE_DATA_READ_URL + '?bit_status=').then((res) => {
            setData(res.data)
        })
    }, [queryType]);

    const renderData = () => {
        switch (queryType){
            case 'Reports':
                break;
            case 'BIT Status':
                break;
        }
        if (data) return <Box>{data}</Box>
        return <div/>;
    }

    return (
    <Box bgcolor="#373A36" minHeight="100vh" textAlign="center">
        <Box bgcolor="#D48166" minHeight="100vh" marginRight="150px" marginLeft="150px">
            <Card>
                <CardHeader title="BIT Framework Query System"/>
                <CardContent>
                    <Grid container justify='space-evenly'>
                        <Selector menuItems={userGroups} currentValue={userGroup} onChange={changeUserGroup} isDisabled={false}/>
                        <Selector menuItems={queryTypes} currentValue={queryType} onChange={changeQueryType} isDisabled={isDisabled}/>
                        <Selector menuItems={filterOptions} currentValue={filter} onChange={changeFilter} isDisabled={isDisabled}/>
                    </Grid>
                    {filter === 'time' &&
                    <Grid className={classes.dateGrid} container justify='space-evenly'>
                        <DatePicker
                            currentDate={startTime}
                            onChange={changeStartTime}
                            placeholder='start time'
                        />
                        <DatePicker
                            currentDate={endTime}
                            onChange={changeEndTime}
                            placeholder={'end time'}
                        />
                    </Grid>}
                </CardContent>
            </Card>
        </Box>
    </Box>
    );
}

export default App;
