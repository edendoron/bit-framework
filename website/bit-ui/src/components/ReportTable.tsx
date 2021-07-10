import React, {FC} from 'react';
import {
    Box,
    Collapse, IconButton,
    makeStyles,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography
} from "@material-ui/core";
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

interface reportObject{
    testId: number,
    reportPriority: number,
    timestamp: Date,
    tagSet: Array<{key: string, value: string}>,
    fieldSet: Array<{key: string, value: string}>,
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

const Row = (props: {row : reportObject}) => {
    const { row } = props;
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    return(
        <React.Fragment>
            <TableRow className={classes.root}>
                <TableCell>
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">{row.testId}</TableCell>
                <TableCell align="center">{row.timestamp}</TableCell>
                <TableCell align="center">{row.reportPriority}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box margin={1}>
                            <Typography variant="h6" gutterBottom component="div">
                                Fields
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        {row.fieldSet.map((field) => (
                                            <TableCell className={classes.tableHeader}>{field.key}</TableCell>
                                        ))}
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow>
                                    {row.fieldSet.map((field) => (
                                            <TableCell component="th" scope="row">
                                                {field.value}
                                            </TableCell>
                                    ))}
                                    </TableRow>
                                </TableBody>
                            </Table>
                            <Typography variant="h6" gutterBottom component="div">
                                Tags
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        {row.tagSet.map((tag) => (
                                            <TableCell className={classes.tableHeader}>{tag.key}</TableCell>
                                        ))}
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow>
                                        {row.tagSet.map((tag) => (
                                            <TableCell component="th" scope="row">
                                                {tag.value}
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                </TableBody>
                            </Table>
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

export const ReportTable: FC<ReportTableProps> = ({data}) => {

    return(
        <TableContainer component={Paper}>
            <Table aria-label="collapsible table">
                <TableHead>
                    <TableRow>
                        <TableCell />
                        <TableCell>Test ID</TableCell>
                        <TableCell align="center">Timestamp</TableCell>
                        <TableCell align="center">Report Priority</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {
                        data.sort(compareDate).map((report) => (
                        <Row key={report.testId} row={report} />
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    )
}