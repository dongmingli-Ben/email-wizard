import React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import CssBaseline from "@mui/material/CssBaseline";
import EastIcon from "@mui/icons-material/East";
import GitHubIcon from "@mui/icons-material/GitHub";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import { useNavigate } from "react-router-dom";
import { IconButton } from "@mui/material";

// TODO remove, this demo shouldn't need to reset the theme.
const defaultTheme = createTheme();

function IntroductionPage() {
  const navigate = useNavigate();
  return (
    <ThemeProvider theme={defaultTheme}>
      <Container
        fixed
        sx={{
          minWidth: "100vw",
          bgcolor: "primary.main",
        }}
      >
        <Container
          component="main"
          maxWidth="sm"
          sx={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            minHeight: "100vh",
          }}
        >
          <CssBaseline />
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              "& .MuiTypography-root": {
                marginBottom: 4,
              },
            }}
          >
            <Typography variant="h2" color="common.white">
              <div className="u-flex u-flex-alignCenter">
                <span>Email Wizard</span>
                <IconButton
                  sx={{ fontSize: "inherit", ml: 2, color: "inherit" }}
                  onClick={() => {
                    window.open(
                      "https://github.com/dongmingli-Ben/email-wizard"
                    );
                  }}
                >
                  <GitHubIcon fontSize="inherit"></GitHubIcon>
                </IconButton>
              </div>
            </Typography>
            <Typography variant="h4" color="common.white">
              Email {<EastIcon></EastIcon>} Calendar
            </Typography>
            {/* <div> */}
            <Typography variant="body1" color="common.white">
              Streamline event management with Email Wizard. Our web app reads
              emails, extracts key details, and populates your calendar
              automatically. Say goodbye to manual entry and inbox clutter. With
              Email Wizard, invitations and activities are effortlessly
              organized. Enjoy seamless integration, as dates, times, locations,
              and participants are recognized instantly. Simplify event planning
              and unlock efficiency with Email Wizard.
            </Typography>
            <Typography variant="subtitle1" color="common.white">
              Preview! Still in progress.
            </Typography>
            <div>
              <Button
                variant="contained"
                color="secondary"
                onClick={() => {
                  navigate("/login");
                }}
              >
                Log In
              </Button>
              <Button
                variant="text"
                color="secondary"
                style={{ marginLeft: "16px" }}
                onClick={() => {
                  navigate("/register");
                }}
              >
                Register
              </Button>
            </div>
          </Box>
        </Container>
      </Container>
    </ThemeProvider>
  );
}

export default IntroductionPage;
