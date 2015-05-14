package ly.stealth.mesos.marathon;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;

public class Util {
    static void copyAndClose(InputStream in, OutputStream out) throws IOException {
        byte[] buffer = new byte[16 * 1024];
        int actuallyRead = 0;

        try(InputStream _in = in; OutputStream _out = out) {
            while (actuallyRead != -1) {
                actuallyRead = in.read(buffer);
                if (actuallyRead != -1) out.write(buffer, 0, actuallyRead);
            }
        }
    }
}
