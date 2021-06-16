import React, {FC} from 'react';
import {createStyles, FormControl, makeStyles, MenuItem, Select} from "@material-ui/core";

interface SelectorProps {
    menuItems: string[],
    queryType: string,
    onChange: (event: React.ChangeEvent<{ value: unknown }>) => void,
}

const useStyles = makeStyles(() =>
    createStyles({
        formControl: {
            maxWidth: 500,
            alignSelf: 'center'
        },
    }),
);

export const Selector: FC<SelectorProps> = ({menuItems, queryType, onChange}) => {
    const classes = useStyles();

    const renderMenuItems = () => {
        return menuItems.map(item => <MenuItem value={item}>{item}</MenuItem>)
    }
    return (
        <FormControl className={classes.formControl}>
            <Select
                value={queryType}
                onChange={onChange}
            >
                {renderMenuItems()}
            </Select>
        </FormControl>
        )
}