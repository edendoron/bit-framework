import React, {FC} from 'react';
import dayjs from "dayjs";
import {
    Box,
    Collapse,
    IconButton,
    makeStyles,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TableSortLabel,
    Typography
} from "@material-ui/core";
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

const dateFormat = 'YYYY-MMMM-DD HH:mm:s';

interface reportObject {
    testId: number,
    reportPriority: number,
    timestamp: Date,
    tagSet: Array<{ key: string, value: string }>,
    fieldSet: Array<{ key: string, value: string }>,
}

interface ReportTableProps {
    data: Array<reportObject>,

}

const useRowStyles = makeStyles({
    root: {
        '& > *': {
            borderBottom: 'set',
            backgroundColor: '#D48166',
        },
    },
    tableHeader: {
        textDecorationLine: 'underline',
        fontWeight: 'bold',
    }
});


const SubTable = (arr: Array<{ key: string, value: string }>) => {
    const classes = useRowStyles();

    return (
        <Table size="small" aria-label="purchases">
            <TableHead>
                <TableRow>
                    {arr.map((item) => (
                        <TableCell className={classes.tableHeader}>
                            {item.key}
                        </TableCell>
                    ))}
                </TableRow>
            </TableHead>
            <TableBody>
                <TableRow>
                    {arr.map((item) => (
                        <TableCell component="th" scope="row">
                            {item.value}
                        </TableCell>
                    ))}
                </TableRow>
            </TableBody>
        </Table>
    )
}

const Row = (props: { row: reportObject }) => {
    const {row} = props;
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    return (
        <React.Fragment>
            <TableRow className={classes.root}>
                <TableCell>
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/>}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row" align="center">{row.testId}</TableCell>
                <TableCell align="center">{dayjs(row.timestamp).format(dateFormat)}</TableCell>
                <TableCell align="center">{row.reportPriority}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{paddingBottom: 0, paddingTop: 0}} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box margin={1}>
                            <Typography variant="h6" gutterBottom component="div">
                                Fields
                            </Typography>
                            {SubTable(row.fieldSet)}
                            <Typography variant="h6" gutterBottom component="div">
                                Tags
                            </Typography>
                            {SubTable(row.tagSet)}
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    )
}

const compareDate = (report1: reportObject, report2: reportObject) => {
    const Date1 = new Date(report1.timestamp).getTime();
    const Date2 = new Date(report2.timestamp).getTime();
    return Date1 > Date2 ? 1 : -1;
}

const comparePriority = (report1: reportObject, report2: reportObject) => {
    return report1.reportPriority > report2.reportPriority ? -1 : 1;
}

const compareID = (report1: reportObject, report2: reportObject) => {
    return report1.testId > report2.testId ? 1 : -1;
}

function stableSort<T>(array: T[], comparator: (a: T, b: T) => number, direction: string) {
    const stabilizedThis = array.map((el, index) => [el, index] as [T, number]);
    stabilizedThis.sort((a, b) => {
        const order = direction == "desc" ? comparator(a[0], b[0]) : -comparator(a[0], b[0]);
        if (order !== 0) return order;
        return a[1] - b[1];
    });
    return stabilizedThis.map((el) => el[0]);
}


const getComparator = (key: string): ((a: reportObject, b: reportObject) => number) => {
    switch (key) {
        case "Date":
            return compareDate
        case "Priority":
            return comparePriority
        case "ID":
            return compareID
        default:
            return compareDate
    }
}

const useHeaderRowStyles = makeStyles({
    root: {
        '& > *': {
            borderBottom: 'set',
            backgroundColor: '#b3a48e',
        },
    },
    text: {
        fontWeight: 'bold',
    }
});

export const ReportTable: FC<ReportTableProps> = ({data}) => {
    const [orderByField, setOrderField] = React.useState("Date")
    const [orderByDirection, setOrderDirection] = React.useState<"asc" | "desc">("asc")

    const classes = useHeaderRowStyles();

    if (data.length === 0) return <div>No Reports Found.</div>

    return (
        <TableContainer component={Paper}>
            <Table aria-label="collapsible table">
                <TableHead className={classes.root}>
                    <TableRow>
                        <TableCell/>
                        <TableCell
                            className={classes.text}
                            align="center"
                            sortDirection={false}>
                            <TableSortLabel
                                active={orderByField === "ID"}
                                direction={orderByDirection}
                                onClick={() => {
                                    if (orderByDirection === "asc")
                                        setOrderDirection("desc")
                                    else
                                        setOrderDirection("asc")
                                    setOrderField("ID")
                                }}>
                                Test ID
                            </TableSortLabel>
                        </TableCell>
                        <TableCell
                            className={classes.text}
                            align="center"
                            sortDirection={false}>
                            <TableSortLabel
                                active={orderByField === "Date"}
                                direction={orderByDirection}
                                onClick={() => {
                                    if (orderByDirection === "asc")
                                        setOrderDirection("desc")
                                    else
                                        setOrderDirection("asc")
                                    setOrderField("Date")
                                }}>
                                Timestamp
                            </TableSortLabel>
                        </TableCell>
                        <TableCell
                            className={classes.text}
                            align="center"
                            sortDirection={false}>
                            <TableSortLabel
                                active={orderByField === "Priority"}
                                direction={orderByDirection}
                                onClick={() => {
                                    if (orderByDirection === "asc")
                                        setOrderDirection("desc")
                                    else
                                        setOrderDirection("asc")
                                    setOrderField("Priority")
                                }}>
                                Report Priority
                            </TableSortLabel>
                        </TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {
                        stableSort(data, getComparator(orderByField), orderByDirection).map((report) => (
                            <Row key={report.testId} row={report}/>

                        ))}
                </TableBody>
            </Table>
        </TableContainer>
    )
}