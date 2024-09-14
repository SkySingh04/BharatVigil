
import argparse
import yaml
import numpy as np
import glob
import os
from PIL import Image
import binascii
from concurrent.futures import ThreadPoolExecutor
from array import *

def load_config(filename):

    try:
        with open(filename, 'r') as file:
            return yaml.safe_load(file)
    except FileNotFoundError:
        print(f"Error: The file '{filename}' was not found.")
        return None
    except yaml.YAMLError as e:
        print(f"Error parsing YAML file: {e}")
        return None
config = load_config('config.yaml')
parser = argparse.ArgumentParser(description='Pass an argument for model inference')
parser.add_argument('input_dir',type=str,help="path of the input_dir")
fifo_path ='/tmp/model_queue'
if not os.path.exists(fifo_path):
    os.mkfifo(fifo_path)

class png_mnsit_creator():
    def __init__(self):
        self.png_saving_dir=config["png_param"]["png_saving_dir"]
        self.png_size=config["png_param"]["png_size"]
        self.pcap_deletion=config["png_param"]["pcap_deletion"]
        self.threads_alloc=config["png_param"]["threads_alloc"]
        self.matrix_size=config["png_param"]["matrix_size"]
        self.mnist_create=config["mnist_param"]["mnist_creation"]
        self.mnist_saving_dir=config["mnist_param"]["mnist_saving_dir"]
        self.png_deletion=config["mnist_param"]["png_deletion"]
        self.trim_flag=config["png_param"]["trim"]
       
    def mnist_creation(self,png_dir):
        data_image= array('B')
        png_f_l=glob.glob(os.path.join(png_dir,"*.png"))
        for iter in png_f_l:
            Im=Image.open(iter)
            pixel=Im.load()
            width, height = Im.size
            # print("original",width,height)
            for x in range(0,width):
                for y in range(0,height):
                    data_image.append(pixel[y,x])
                    #print([y,x])
            if self.png_deletion:
                os.remove(iter)

        hexval = "{0:#0{1}x}".format(len(png_f_l),6) # number of files in HEX
        hexval = '0x' + hexval[2:].zfill(8)  
        header = array('B')
        header.extend([0,0,8,1])
        header.append(int('0x'+hexval[2:][0:2],16))
        header.append(int('0x'+hexval[2:][2:4],16))
        header.append(int('0x'+hexval[2:][4:6],16))
        header.append(int('0x'+hexval[2:][6:8],16))
        width, height = Im.size
        if max([width,height]) <= 256:
            header.extend([0,0,0,width,0,0,0,height])
        else:
            raise ValueError('Image exceeds maximum size: 256x256 pixels')
        output_file=open(f"{self.mnist_saving_dir}/{self.folder_name}-images-idx3-ubyte",'wb')
        header[3] = 3 # Changing MSB for image data (0x00000803)
        data_image = header + data_image
        data_image.tofile(output_file)

    def trim(self,file):
        # file_len=len(file)
        file_len = len(file)
        
        len_diff = file_len - self.matrix_size
        
        if len_diff > 0:
            # Trim the content
            file = file[:self.matrix_size]
            
        elif len_diff < 0:
            # Pad the content with zero bytes
            
            padding_len = -len_diff
            
            file = file+bytes(padding_len)  # yaha pe error aara hai
            # print(file)
        
        return file
    
    def getMatrixfrom_pcap(self,filename,width):
        with open(filename, 'rb') as f:
            content = f.read()
            
        if self.trim_flag:
            content=self.trim(content) 

       
        hexst = binascii.hexlify(content)  
        fh = np.array([int(hexst[i:i+2],16) for i in range(0, len(hexst), 2)]) 
               
        rn = int(len(fh)/width)
       
        fh = np.reshape(fh[:rn*width],(-1,width))  
        fh = np.uint8(fh)
        
        return fh



    def png_creation(self,file):
        im=Image.fromarray(self.getMatrixfrom_pcap(file,self.png_size))
        
        dir_name=f"{self.png_saving_dir}/{self.folder_name}"
        
        os.makedirs(dir_name,exist_ok=True)
        filename=os.path.basename(file)
        saving_loc=f"{self.png_saving_dir}/{self.folder_name}/{filename}"+".png"
        im.save(saving_loc)
        with open(fifo_path,'w') as fifo:
            fifo.write(f"{self.png_saving_dir}/{self.folder_name}"+'\n')
            fifo.flush()
        if self.pcap_deletion:
            os.remove(file)


    def png_creation_dir(self,input_dir):
        pcap_file_l=glob.glob(os.path.join(input_dir,"*.pcap"))
        self.folder_name=os.path.basename(os.path.normpath(input_dir))
        with ThreadPoolExecutor(max_workers=self.threads_alloc) as executor:
            executor.map(self.png_creation,pcap_file_l)
        if self.mnist_create:
            self.mnist_creation(f"{self.png_saving_dir}/{self.folder_name}")
        
if __name__=="__main__":
    args=parser.parse_args()
    creator=png_mnsit_creator()
    creator.png_creation_dir(args.input_dir)
