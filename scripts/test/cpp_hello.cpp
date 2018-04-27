#include <iostream>
#include <unistd.h>
using namespace std;

int main() {
    for(int i=0; i<10000; i++) {
//        cout<<"hello world"<<endl;
        int a = i;
    }
    sleep(5);
    cout<<"end"<<endl;
    return 0;
}