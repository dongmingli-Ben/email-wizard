import React from "react";

import "./UserAccountInfo.css";
import { Box, Grid, IconButton, Typography } from "@mui/material";
import AddCircleOutlineIcon from "@mui/icons-material/AddCircleOutline";

type UserAccountInfoProps = {
  userName: string;
  userAccounts: { address: string; protocol: string }[];
  setAddAccount: (status: boolean) => void;
};

const UserAccountInfo = (props: UserAccountInfoProps) => {
  let rows = props.userAccounts.map(
    (account: { address: string; protocol: string }, index: number) => {
      return (
        <Grid
          item
          sx={{
            width: "100%",
          }}
          key={index}
        >
          <Box
            sx={{
              "&:hover": {
                backgroundColor: "primary.dark",
                opacity: [0.9, 0.8, 0.7],
                cursor: "default",
              },
              width: "100%",
              pl: "10%",
              pr: "5%",
              pt: 1,
              pb: 1,
              boxSizing: "border-box",
            }}
          >
            <Typography
              color="inherit"
              sx={{
                width: "100%",
              }}
              noWrap
            >
              {account.address}
            </Typography>
          </Box>
        </Grid>
      );
    }
  );
  return (
    <Grid
      container
      direction="column"
      justifyContent="flex-start"
      sx={{
        color: "common.white",
        mt: 3,
        width: "100%",
      }}
    >
      <Grid
        item
        sx={{
          width: "100%",
        }}
      >
        <Box
          sx={{
            "&:hover": {
              backgroundColor: "primary.dark",
              opacity: [0.9, 0.8, 0.7],
              cursor: "default",
            },
            width: "100%",
            pl: "10%",
            pr: "5%",
            pt: 1,
            pb: 1,
            boxSizing: "border-box",
          }}
        >
          <Typography
            color="inherit"
            variant="h4"
            noWrap
            sx={{
              width: "100%",
            }}
          >
            {props.userName}
          </Typography>
        </Box>
      </Grid>
      {props.userAccounts.length > 0 ? rows : <></>}
      <Grid item>
        <Box
          sx={{
            "&:hover": {
              backgroundColor: "primary.dark",
              opacity: [0.9, 0.8, 0.7],
              cursor: "default",
            },
            width: "100%",
            pl: "10%",
            pr: "5%",
            boxSizing: "border-box",
          }}
        >
          <IconButton
            onClick={() => {
              props.setAddAccount(true);
            }}
            sx={{
              color: "inherit",
              pl: 0,
              pb: 0,
            }}
          >
            <div className="u-flex u-flex-alignCenter">
              <span>
                <AddCircleOutlineIcon color="inherit"></AddCircleOutlineIcon>
              </span>
              <span>
                <Typography
                  color="inherit"
                  sx={{
                    width: "100%",
                    fontWeight: "bold",
                    ml: 0.5,
                  }}
                  noWrap
                  gutterBottom
                >
                  Add Mailbox
                </Typography>
              </span>
            </div>
          </IconButton>
        </Box>
      </Grid>
    </Grid>
  );
};

export default UserAccountInfo;
