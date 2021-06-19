import React, {useCallback, useEffect, useState} from 'react';
import './App.css';
import {getReports, getBitStatus} from './utils/queryAPI';
import {Box, Button, Card, CardContent, CardHeader, createStyles, Grid, makeStyles, Paper} from '@material-ui/core'
import {Selector} from "./components/Selector";
import {DatePicker} from "./components/DatePicker";

const queryTypes = ['Reports', 'BIT Status'];
const userGroups = ['group1', 'group2', 'group3', 'group4', 'groupRafael', 'TemperatureCelsius group', 'group general', 'groupField'];
const filterOptions = ['time', 'tag', 'field'];

const useStyles = makeStyles(() =>
    createStyles({
        dateGrid: {
            marginTop: 20,
        },
        paper: {
            width: '50%',
            marginLeft: '25%',
            marginTop: 20,
        }
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
    const [data, setData] = useState('');

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

    const renderData = async () => {
        switch (queryType){
            case 'Reports':
                setData(await getReports(filter, startTime, endTime));
                break;
            case 'BIT Status':
                setData(await getBitStatus(userGroup, startTime, endTime, filter));
                break;
        }
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
                <Button onClick={renderData}>
                    send
                </Button>
            </Card>
            <Paper className={classes.paper}>{JSON.stringify(data)}</Paper>
        </Box>
    </Box>
    );
}

export default App;
