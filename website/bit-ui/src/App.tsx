import React, {useEffect, useState} from 'react';
import './App.css';
import {getBitStatus, getReports, getUserGroups} from './utils/queryAPI';
import {
    Box,
    Button,
    Card,
    CardContent,
    CardHeader,
    createMuiTheme,
    createStyles,
    Grid,
    makeStyles,
    MuiThemeProvider, TextField,
} from '@material-ui/core'
import {Selector} from "./components/Selector";
import {DatePicker} from "./components/DatePicker";
import {ReportTable} from "./components/ReportTable";
import {StatusTable} from "./components/StatusTable";


const queryTypes = ['Reports', 'BIT Status'];
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
            },
            sendButton: {
                border: '3px solid #D48166 ',
                marginBottom: 10,
            },
            textInput: {
                marginLeft: 20,
                marginTop: 20
            }
        },
    ));


export const App = () => {
    const classes = useStyles();

    const [queryType, setQueryType] = useState('');
    const [userGroup, setUserGroup] = useState('');
    const [userGroups, setUserGroups] = useState(['']);
    const [filter, setFilter] = useState('');
    const [isDisabled, setDisabled] = useState(true);
    const [startTime, setStartTime] = useState(new Date());
    const [endTime, setEndTime] = useState(new Date());
    const [data, setData] = useState<string>();
    const [error, setError] = useState(null);
    const [field, setField] = useState('');
    const [tagKey, setTagKey] = useState('');
    const [tagValue, setTagValue] = useState('');

    useEffect(() => {
        async function fetchUserGroups() {
            setUserGroups(await getUserGroups())
        }
        fetchUserGroups()
    }, [])

    const changeQueryType = (event: React.ChangeEvent<{ value: unknown }>) => {
        setQueryType(event.target.value as string);
        setFilter('');
        setData('');
    }

    const changeUserGroup = async (event: React.ChangeEvent<{ value: unknown }>) => {
        setUserGroup(event.target.value as string);
        setFilter('');
        setQueryType('');
        setDisabled(false);
    }

    const changeFilter = (event: React.ChangeEvent<{ value: unknown }>) => {
        setFilter((event.target.value as string));
        setTagValue('');
        setTagKey('');
        setField('');
    }

    const changeStartTime = (date: Date) => {
        setStartTime(date)
    }

    const changeEndTime = (date: Date) => {
        setEndTime(date)
    }

    const getData = async () => {
        switch (queryType) {
            case 'Reports':
                try {
                    setData(await getReports(filter, startTime, endTime, field, tagKey, tagValue));
                    setError(null);
                } catch (e) {
                    setError(e);
                }
                break;
            case 'BIT Status':
                try {
                    setData(await getBitStatus(userGroup, startTime, endTime, filter));
                    setError(null);
                } catch (e) {
                    setError(e);
                }
                break;
        }
    }

    const renderData = () => {
        if (error) return <div>Network error. Please try again.</div>
        switch (queryType) {
            case 'Reports':
                if (!!data) {
                    return <ReportTable data={JSON.parse(data)}/>
                }
                break;
            case 'BIT Status':
                if (!!data) {
                    return <StatusTable data={JSON.parse(data)}/>
                }
                break;
        }
    }

    const theme = createMuiTheme({
        palette: {
            primary: {
                main: '#D48166'
            },
        }
    })

    return (
        <MuiThemeProvider theme={theme}>
            <Box bgcolor="#373A36" minHeight="100vh" textAlign="center">
                <Box bgcolor="#D48166" minHeight="100vh" marginRight="150px" marginLeft="150px">
                    <Card>
                        <CardHeader title="BIT Framework Query System"/>
                        <CardContent>
                            <Grid container justify='space-evenly'>
                                <Selector menuItems={userGroups} currentValue={userGroup} onChange={changeUserGroup}
                                          isDisabled={false} placeholder='user group'/>
                                <Selector menuItems={queryTypes} currentValue={queryType} onChange={changeQueryType}
                                          isDisabled={isDisabled} placeholder='query type'/>
                                <Selector menuItems={filterOptions} currentValue={filter} onChange={changeFilter}
                                          isDisabled={isDisabled} placeholder='filter'/>
                            </Grid>
                            <Grid className={classes.dateGrid} container justify='space-evenly'>
                                <DatePicker
                                    currentDate={startTime}
                                    onDateChange={changeStartTime}
                                    label='start time'
                                />
                                <DatePicker
                                    currentDate={endTime}
                                    onDateChange={changeEndTime}
                                    label='end time'
                                />
                            </Grid>
                            {filter === 'field' &&
                                <TextField
                                    label='field name'
                                    onChange={(e) => setField(e.target.value)}
                                    className={classes.textInput}
                                    variant='filled'
                                />
                            }
                            {filter === 'tag' &&
                                <>
                                    <TextField
                                        label='tag key'
                                        onChange={(e) => setTagKey(e.target.value)}
                                        className={classes.textInput}
                                        variant='filled'
                                    />
                                    <TextField
                                        label='tag value'
                                        onChange={(e) => setTagValue(e.target.value)}
                                        className={classes.textInput}
                                        variant='filled'
                                    />
                                </>
                            }
                        </CardContent>
                        <Button className={classes.sendButton} onClick={getData} disabled={queryType === '' || filter === ''}>
                            Send
                        </Button>
                    </Card>
                    {renderData()}
                </Box>
            </Box>
        </MuiThemeProvider>
    );
}

export default App;
