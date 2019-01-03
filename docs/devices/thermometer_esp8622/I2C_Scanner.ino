/*
----------------------------------------------------------------------
I2C-Scanner with text diagnsotics only

Code fuer Arduino
Author Retian
Version 1.0

Der Sketch testet die Standard-7-Bit-Adressen von 1 bis 128
mit den Frequenzen 100 kHz und 400 kHz

ToDo:
- add scanblock diagram - separate version

more Infos: https://learn.adafruit.com/i2c-addresses/the-list
more Infos: https://arduino-projekte.webnode.at/i2c-scanner/

Andreas Meier - 11.11.2018

----------------------------------------------------------------------
*/

#include <Wire.h>

void I2C_Scanner()
{
  Serial.println("\n2C-Scanner with text diagnsotics only");

  byte error, address;
  int nDevices;

  Serial.println("Scanning ...");

  nDevices = 0;
  for(address = 1; address < 128; address++ )
  {
    // The i2c_scanner uses the return value of
    // the Write.endTransmisstion to see if
    // a device did acknowledge to the address.
    Wire.beginTransmission(address);
    error = Wire.endTransmission();

    if (error == 0)
    {
      Serial.print("I2C device found at address 0x");
      if (address<16)
        Serial.print("0");
      Serial.print(address,HEX);
      Serial.println("  !");

      nDevices++;
    }
    else if (error==4)
    {
      Serial.print("Unknow error at address 0x");
      if (address<16)
        Serial.print("0");
      Serial.println(address,HEX);
    }

  } // end for

  // end for loop
  Serial.println();

  if (nDevices == 0)
    while(1)
    {
      Serial.println("No I2C devices found !  ... please Restart ");
      delay(1000);
    }
  else
    Serial.println("done\n");

}


// -----------------------------------------------------------------------
